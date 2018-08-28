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
		Addr:           config.Port,
		ReadTimeout:    60 * time.Second,
		WriteTimeout:   60 * time.Second,
		MaxHeaderBytes: 1 << 16,
	}

	// Setup providers and endpoints for OAuth
	userStore := NewUserStore(db.DB)
	authHandler := auth.New(userStore)
	authHandler.SetCORS(config.DevURL)

	for name, provider := range config.Providers {
		callbackURL, err := url.Parse(config.URL)
		if config.DevURL != "" {
			callbackURL, err = url.Parse(config.DevURL)
		}
		if err != nil {
			panic(err)
		}
		callbackURL.Path = path.Join(callbackURL.Path, provider.URI)

		log.Println("API callback:", name, callbackURL.String())
		authHandler.AddProvider(name, provider.ID, provider.Secret, callbackURL.String(), provider.Scopes)
	}

	http.HandleFunc("/auth/list", authHandler.Auth)
	http.HandleFunc("/auth/token", authHandler.Token)

	// Endpoints for the API and Vue client
	vueHandler := http.FileServer(Vue("client/dist/"))
	apiHandler := api.New(db)
	apiHandler.SetCORS(config.DevURL)

	http.Handle("/api/", authHandler.Middleware(apiHandler))
	http.Handle("/", vueHandler)

	log.Println("Listening on", config.Port, "at", config.URL)
	log.Fatal(s.ListenAndServe())
}

type ConfigProvider struct {
	ID     string   `json:"id"`
	Secret string   `json:"secret"`
	URI    string   `json:"uri"`
	Scopes []string `json:"scopes"`
}

type Config struct {
	Name      string
	Port      string
	URL       string
	DevURL    string
	Database  string
	Scheme    string
	Providers map[string]ConfigProvider
}

func loadConfig() *Config {
	configFilename := ""

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options] [input]\n\nOptions:\n", os.Args[0])
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nInput:\n  Files or directories, leave blank to use stdin\n")
	}
	flag.StringVarP(&configFilename, "config", "", "config.json", "Config filename")
	flag.Parse()

	configFile, err := os.Open(configFilename)
	if err != nil {
		panic(err)
	}
	defer configFile.Close()

	config := &Config{
		Port:     ":3000",
		URL:      "http://localhost:8080",
		Database: "sqlite.db",
		Scheme:   "scheme.sql",
	}
	if err := json.NewDecoder(configFile).Decode(config); err != nil {
		panic("parsing config: " + err.Error())
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
