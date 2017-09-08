package api

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gorilla/context"
)

type API struct {
	db      *sql.DB
	corsURL string
}

func New(db *sql.DB) *API {
	return &API{
		db,
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

	user := context.Get(r, "user")
	if user == nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	fmt.Fprintln(w, "API endpoint for", user)
}
