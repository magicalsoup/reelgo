package users

import (
	"database/sql"

	"github.com/gorilla/mux"
)


func AddRoutes(r *mux.Router, db* sql.DB) {
	r.HandleFunc("/login", loginHandler(db)).Methods("POST")
	r.HandleFunc("/signup", signUpHandler(db)).Methods("POST", "OPTIONS")
}