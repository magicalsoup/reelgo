package users

import (
	"database/sql"

	"github.com/gorilla/mux"
)

func AddRoutes(r *mux.Router, db* sql.DB) {
	r.HandleFunc("/login", loginHandler(db)).Methods("GET")
	r.HandleFunc("/signup", signUpHandler(db)).Methods("POST")
}