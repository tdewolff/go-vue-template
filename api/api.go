package api

import (
	"database/sql"
	"fmt"
	"net/http"
)

type API struct {
	db *sql.DB
}

func New(db *sql.DB) *API {
	return &API{
		db,
	}
}

func (api *API) CORS(clientURL string, next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", clientURL)
		w.Header().Set("Access-Control-Allow-Methods", "GET,PUT,POST,DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type,Authorization")

		next(w, r)
	})
}

func (api *API) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == "OPTIONS" {
		return
	}

	fmt.Fprintln(w, "API endpoint")
}
