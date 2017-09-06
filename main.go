package main

import (
	"database/sql"
	"log"
	"net/http"
	"path"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/tdewolff/go-vue-template/auth"
)

const clientURL = "http://localhost:8080"
const authCallbackURI = "/auth/callback"
const jwtSecret = "RANDOMSECRET" // Change this!
const jwtDuration = time.Hour * 24

func main() {
	db, err := sql.Open("sqlite3", "./sqlite.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	s := &http.Server{
		Addr:           ":3000",
		ReadTimeout:    60 * time.Second,
		WriteTimeout:   60 * time.Second,
		MaxHeaderBytes: 1 << 16,
	}

	a := auth.New(clientURL, authCallbackURI, jwtSecret, jwtDuration)
	a.AddProvider(auth.Google("key", "secret"))
	a.AddProvider(auth.Facebook("key", "secret"))
	a.AddProvider(auth.GitHub("key", "secret"))
	http.HandleFunc("/auth/list", a.List)
	http.HandleFunc("/auth/jwt", a.JWT)

	http.Handle("/", http.FileServer(Vue("client/dist/")))

	log.Fatal(s.ListenAndServe())
}

type Vue string

func (v Vue) Open(name string) (http.File, error) {
	if ext := path.Ext(name); name != "/" && (ext == "" || ext == ".html") {
		name = "index.html"
	}
	return http.Dir(v).Open(name)
}
