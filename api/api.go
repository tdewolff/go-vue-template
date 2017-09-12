package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/context"
	"github.com/tdewolff/auth"
	"golang.org/x/oauth2"
)

type API struct {
	db      *sql.DB
	auth    *auth.Auth
	corsURL string
}

func New(db *sql.DB, auth *auth.Auth) *API {
	return &API{
		db,
		auth,
		"",
	}
}

func (api *API) SetCORS(clientURL string) {
	api.corsURL = clientURL
}

func (api *API) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if api.corsURL != "" {
		w.Header().Set("Access-Control-Allow-Origin", api.corsURL)
		w.Header().Set("Access-Control-Allow-Methods", "GET,PUT,POST,DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type,Authorization")
	}

	if r.Method == "OPTIONS" {
		return
	}

	user := context.Get(r, "user").(*auth.User)
	if user == nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	fmt.Fprintln(w, "API endpoint for", user)
}

func (api *API) getClient(user *auth.User, providerID string) *http.Client {
	var tokenString string
	if err := api.db.QueryRow(`SELECT token FROM social_tokens WHERE user_id=? AND provider=?`, user.ID, providerID).Scan(&tokenString); err != nil {
		if err != sql.ErrNoRows {
			log.Println("cannot get token for %v: %v\n", providerID, err)
		}
		return nil
	}

	var token *oauth2.Token
	if err := json.Unmarshal([]byte(tokenString), &token); err != nil {
		log.Println("cannot decode token for %v: %v\n", providerID, err)
		return nil
	}

	provider := api.auth.GetProvider(providerID)
	if provider == nil {
		log.Println("cannot get provider for %v\n", providerID)
		return nil
	}
	return provider.Config.Client(oauth2.NoContext, token)
}
