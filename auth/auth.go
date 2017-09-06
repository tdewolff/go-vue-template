package auth

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/securecookie"
)

var sessionDuration = time.Minute * 10

type Profile struct {
	Name  string
	Email string
}

type Session struct {
	csrf       string
	refererURI string
	expires    time.Time
}

type Auth struct {
	clientURL   string
	redirectURL string
	jwtSecret   string
	jwtDuration time.Duration
	providers   []Provider

	sec      *securecookie.SecureCookie
	sessions map[string]Session
}

func New(clientURL, redirectURI, jwtSecret string, jwtDuration time.Duration) *Auth {
	return &Auth{
		clientURL,
		clientURL + redirectURI,
		jwtSecret,
		jwtDuration,
		[]Provider{},

		securecookie.New(GenerateSecret(32), nil),
		map[string]Session{},
	}
}

func (a *Auth) AddProvider(provider Provider) {
	provider.SetRedirectURL(a.redirectURL)
	a.providers = append(a.providers, provider)
}

type ProviderList struct {
	SessionID string         `json:"sessionId"`
	Providers []ProviderItem `json:"providers"`
}

// ProviderItem is the response for the List request
type ProviderItem struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	URL  string `json:"url`
}

func (a *Auth) List(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", a.clientURL)
	w.Header().Set("Access-Control-Allow-Methods", "GET,PUT,POST,DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// Early return for pre-flighted requests, deny other requests except GET and POST
	if r.Method == "OPTIONS" {
		return
	} else if r.Method != "GET" && r.Method != "POST" {
		http.Error(w, "403 forbidden", http.StatusForbidden)
		return
	}

	// Already authorized clients are forbidden
	if authorization := r.Header.Get("Authorization"); authorization != "" {
		http.Error(w, "403 forbidden", http.StatusForbidden)
		return
	}

	// Get referer URI from parameters
	r.ParseForm()
	refererURI := r.Form.Get("referer_uri")

	// Generate Cross-Site Request Forgery token and store in memory
	csrf := base64.URLEncoding.EncodeToString(GenerateSecret(32))
	sessionID := base64.URLEncoding.EncodeToString(GenerateSecret(32))
	a.sessions[sessionID] = Session{
		csrf,
		refererURI,
		time.Now().Add(sessionDuration),
	}

	// Encode session ID
	encodedSessionID, err := a.sec.Encode("auth", sessionID)
	if err != nil {
		log.Println("could not encode session ID:", err)
		http.Error(w, "500 internal server error", http.StatusInternalServerError)
		return
	}

	// Generate states for every provider, save them to the session store and generate an authorization URL
	list := ProviderList{encodedSessionID, make([]ProviderItem, 0, len(a.providers))}
	for _, provider := range a.providers {
		state := encodeState(csrf, provider.ID(), refererURI)
		item := ProviderItem{
			provider.ID(),
			provider.Name(),
			provider.AuthURL(state),
		}
		list.Providers = append(list.Providers, item)
	}

	// Send JSON response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(list); err != nil {
		log.Println("could not encode list response:", err)
		http.Error(w, "500 internal server error", http.StatusInternalServerError)
		return
	}
}

func (a *Auth) JWT(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", a.clientURL)
	w.Header().Set("Access-Control-Allow-Methods", "GET,PUT,POST,DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// Early return for pre-flighted requests, deny other requests except GET and POST
	if r.Method == "OPTIONS" {
		return
	} else if r.Method != "GET" && r.Method != "POST" {
		http.Error(w, "403 forbidden", http.StatusForbidden)
		return
	}

	// Already authorized clients are forbidden
	if authorization := r.Header.Get("Authorization"); authorization != "" {
		http.Error(w, "403 forbidden", http.StatusForbidden)
		return
	}

	// Get code and state parameters, and unpack state
	r.ParseForm()
	code := r.Form.Get("code")
	encodedSessionID := r.Form.Get("session_id")
	csrf, providerID, refererURI, err := decodeState(r.Form.Get("state"))
	if err != nil {
		log.Println("Decoding state failed:", err)
		http.Error(w, "500 internal server error", http.StatusInternalServerError)
		return
	}

	// Decode session ID
	sessionID := ""
	if err := a.sec.Decode("auth", encodedSessionID, &sessionID); err != nil {
		log.Println("could not decode session ID:", err)
		http.Error(w, "500 internal server error", http.StatusInternalServerError)
		return
	}

	// Remove expired sessions
	for id, session := range a.sessions {
		if time.Now().After(session.expires) {
			fmt.Println("delete session", id)
			delete(a.sessions, id)
		}
	}

	// Check that state equals previously generated state, to ensure that the client requested access through this server
	if session, ok := a.sessions[sessionID]; ok && (csrf != session.csrf || refererURI != session.refererURI) {
		log.Printf("Bad state for %s: ok %v and csrf %s != %s or refererURI %s != %s\n", providerID, ok, csrf, session.csrf, refererURI, session.refererURI)
		http.Error(w, "403 forbidden", http.StatusForbidden)
		return
	}

	// Search for provider
	var provider Provider
	for _, p := range a.providers {
		if p.ID() == providerID {
			provider = p
			break
		}
	}
	if provider == nil {
		log.Println("Provider does not exist:", providerID)
		http.Error(w, "403 forbidden", http.StatusForbidden)
		return
	}

	// Get profile data from the provider
	profile, err := provider.Profile(code)
	if err != nil {
		log.Printf("OAuth request to %v failed: %v\n", provider.Name(), err)
		http.Error(w, "403 forbidden", http.StatusForbidden)
		return
	}

	// Create and sign JWT
	jwt := generateJWT(a.jwtSecret, a.jwtDuration, profile)

	// Send JSON response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	v := struct {
		JWT string `json:"jwt"`
		URI string `json:"uri"`
	}{
		jwt,
		refererURI,
	}
	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Println("could not encode auth response:", err)
		http.Error(w, "500 internal server error", http.StatusInternalServerError)
		return
	}
}

// GenerateSecret returns a byte slice of length n of cryptographically secure random data
func GenerateSecret(n int) []byte {
	secret := make([]byte, n)
	if _, err := rand.Read(secret); err != nil {
		panic(err)
	}
	return secret
}

func encodeState(csrf, providerID, refererURI string) string {
	values := make(url.Values, 3)
	values.Add("sec", csrf)
	values.Add("prv", providerID)
	values.Add("uri", refererURI)
	return base64.URLEncoding.EncodeToString([]byte(values.Encode()))
}

func decodeState(state string) (string, string, string, error) {
	query, err := base64.URLEncoding.DecodeString(state)
	if err != nil {
		return "", "", "", err
	}
	values, err := url.ParseQuery(string(query))
	if err != nil {
		return "", "", "", err
	}
	return values.Get("sec"), values.Get("prv"), values.Get("uri"), nil
}

func generateJWT(jwtSecret string, jwtDuration time.Duration, profile Profile) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp":   time.Now().Add(jwtDuration).Unix(),
		"name":  profile.Name,
		"email": profile.Email,
	})
	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		panic(err)
	}
	return tokenString
}
