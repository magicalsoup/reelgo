package auth

import (
	"database/sql"

	"github.com/gorilla/mux"
)

func AddRoutes(r *mux.Router, db* sql.DB) {
	r.HandleFunc("/auth", codeCheckHandler(db)).Methods("POST")
}