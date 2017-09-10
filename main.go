package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"time"

	_ "github.com/mattn/go-sqlite3"
	flag "github.com/spf13/pflag"
	"github.com/tdewolff/auth"
	"github.com/tdewolff/go-vue-template/api"
	"github.com/tdewolff/minify"
	"github.com/tdewolff/minify/html"
)

var urlScheme = "http"
var jwtSecret = []byte("RANDOMSECRET") // Change this!
var jwtDuration = time.Minute * 30

var db *sql.DB

func main() {
	config := loadConfig()

	db = initDatabase(config)
	defer db.Close()

	s := &http.Server{
		Addr:           ":" + config.port,
		ReadTimeout:    60 * time.Second,
		WriteTimeout:   60 * time.Second,
		MaxHeaderBytes: 1 << 16,
	}

	// Setup providers and endpoints for OAuth
	authHandler := auth.New(config.Name, &UserStore{}, jwtSecret, jwtDuration)
	defer authHandler.Close()

	if config.dev {
		authHandler.SetCORS(config.clientURL)
	}

	for id, provider := range config.Providers {
		callbackURL, err := url.Parse(config.clientURL)
		if err != nil {
			panic(err)
		}
		callbackURL.Path = path.Join(callbackURL.Path, provider.URI)

		authHandler.AddProvider(id, provider.ID, provider.Secret, callbackURL.String())
	}

	http.HandleFunc("/auth/list", authHandler.AuthList)
	http.HandleFunc("/auth/token", authHandler.Token)

	// Endpoints for the API and Vue client
	vueHandler := http.FileServer(Vue("client/dist/"))
	apiHandler := api.New(db)
	if config.dev {
		apiHandler.SetCORS(config.clientURL)
	}

	m := minify.New()
	m.AddFunc("text/html", html.Minify)

	http.Handle("/api/", authHandler.Middleware(apiHandler))
	http.Handle("/", m.Middleware(vueHandler))

	log.Fatal(s.ListenAndServe())
}

type ConfigProvider struct {
	ID     string `json:"id"`
	Secret string `json:"secret"`
	URI    string `json:"uri"`
}

type Config struct {
	Name      string                    `json:"name"`
	Host      string                    `json:"host"`
	HostDev   string                    `json:"hostDev"`
	Database  string                    `json:"database"`
	Scheme    string                    `json:"scheme"`
	Providers map[string]ConfigProvider `json:"providers"`

	dev       bool
	port      string
	clientURL string
}

func loadConfig() *Config {
	configFilename := ""
	dev := false

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options] [input]\n\nOptions:\n", os.Args[0])
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nInput:\n  Files or directories, leave blank to use stdin\n")
	}
	flag.StringVarP(&configFilename, "config", "", "config.json", "Config filename")
	flag.BoolVarP(&dev, "dev", "", false, "Development mode enables HostDev for client URLs and sets CORS headers")
	flag.Parse()

	configFile, err := os.Open(configFilename)
	if err != nil {
		panic(err)
	}
	defer configFile.Close()

	config := &Config{
		Host:     "localhost:3000",
		HostDev:  "localhost:8080",
		Database: "sqlite.db",
		Scheme:   "scheme.sql",
		dev:      dev,
	}
	if err := json.NewDecoder(configFile).Decode(config); err != nil {
		panic("parsing config: " + err.Error())
	}

	serverURL, err := url.Parse(urlScheme + "://" + config.Host)
	if err != nil {
		panic(err)
	}
	config.port = serverURL.Port()
	config.clientURL = serverURL.String()
	if dev {
		clientURL, err := url.Parse(urlScheme + "://" + config.HostDev)
		if err != nil {
			panic(err)
		}
		config.clientURL = clientURL.String()
	}

	return config
}

func initDatabase(config *Config) *sql.DB {
	db, err := sql.Open("sqlite3", config.Database)
	if err != nil {
		panic(err)
	}

	schemeFile, err := os.Open(config.Scheme)
	if err != nil {
		panic(err)
	}
	defer schemeFile.Close()

	scheme, err := ioutil.ReadAll(schemeFile)
	if err != nil {
		panic(err)
	}
	if _, err = db.Exec(string(scheme)); err != nil {
		panic("parsing scheme: " + err.Error())
	}
	return db
}

////////////////

type UserStore struct{}

// Login logs in every user, creating a new account if it is a new user
func (*UserStore) Login(user *auth.User, provider, accessToken, refreshToken string) (bool, error) {
	if err := db.QueryRow(`SELECT id, name FROM users WHERE email=?`, user.Email).Scan(&user.ID, &user.Name); err != nil && err != sql.ErrNoRows {
		return false, err
	} else if err == sql.ErrNoRows {
		res, err := db.Exec(`INSERT INTO users (name, email) VALUES (?, ?)`, user.Name, user.Email)
		if err != nil {
			return false, err
		}

		user.New = true
		user.ID, err = res.LastInsertId()
		if err != nil {
			return false, err
		}
	}

	_, err := db.Exec(`INSERT OR REPLACE INTO social_tokens (user_id, provider, access_token, refresh_token) VALUES (?, ?, ?, ?)`, user.ID, provider, accessToken, refreshToken)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (*UserStore) Validate(user *auth.User) (bool, error) {
	var id int
	if err := db.QueryRow(`SELECT id FROM users WHERE id=? AND email=?`, user.ID, user.Email).Scan(&id); err != nil && err != sql.ErrNoRows {
		return false, err
	} else if err == sql.ErrNoRows {
		return false, nil
	}
	return true, nil
}

////////////////

type Vue string

func (v Vue) Open(name string) (http.File, error) {
	if ext := path.Ext(name); name != "/" && (ext == "" || ext == ".html") {
		name = "index.html"
	}
	return http.Dir(v).Open(name)
}
