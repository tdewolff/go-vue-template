package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"time"

	"github.com/gorilla/sessions"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	flag "github.com/spf13/pflag"
	"github.com/tdewolff/auth"
	"github.com/tdewolff/go-vue-template/api"
)

const urlScheme = "http"

func main() {
	config := loadConfig()
	db := initDatabase(config)

	s := &http.Server{
		Addr:           ":" + config.port,
		ReadTimeout:    60 * time.Second,
		WriteTimeout:   60 * time.Second,
		MaxHeaderBytes: 1 << 16,
	}

	// Setup providers and endpoints for OAuth
	sessionStore := sessions.NewCookieStore([]byte("something-very-secret"))
	authHandler := auth.New(sessionStore, auth.NewDefaultUserStore(db.DB))

	if config.dev {
		authHandler.SetDevURL(config.clientURL)
	}

	for name, provider := range config.Providers {
		callbackURL, err := url.Parse(config.clientURL)
		if err != nil {
			panic(err)
		}
		callbackURL.Path = path.Join(callbackURL.Path, provider.URI)

		authHandler.AddProvider(name, provider.ID, provider.Secret, callbackURL.String(), provider.Scopes)
	}

	http.HandleFunc("/auth/list", authHandler.Auth)
	http.HandleFunc("/auth/token", authHandler.Token)
	http.HandleFunc("/auth/logout", authHandler.Logout)

	// Endpoints for the API and Vue client
	vueHandler := http.FileServer(Vue("client/dist/"))
	apiHandler := api.New(db)
	if config.dev {
		apiHandler.SetDevURL(config.clientURL)
	}

	http.Handle("/api/", authHandler.Middleware(apiHandler))
	http.Handle("/", vueHandler)

	log.Fatal(s.ListenAndServe())
}

type ConfigProvider struct {
	ID     string   `json:"id"`
	Secret string   `json:"secret"`
	URI    string   `json:"uri"`
	Scopes []string `json:"scopes"`
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

func initDatabase(config *Config) *sqlx.DB {
	db := sqlx.MustConnect("sqlite3", config.Database)

	schemeFile, err := os.Open(config.Scheme)
	if err != nil {
		panic(err)
	}
	defer schemeFile.Close()

	scheme, err := ioutil.ReadAll(schemeFile)
	if err != nil {
		panic(err)
	}

	db.MustExec(string(scheme))
	return db
}

////////////////

type Vue string

func (v Vue) Open(name string) (http.File, error) {
	if ext := path.Ext(name); name != "/" && (ext == "" || ext == ".html") {
		name = "index.html"
	}
	return http.Dir(v).Open(name)
}
