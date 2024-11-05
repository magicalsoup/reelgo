package users

import (
	"database/sql"

	"github.com/gorilla/mux"
)


func AddRoutes(r *mux.Router, db* sql.DB) {
	r.HandleFunc("/login", loginHandler(db))
	r.HandleFunc("/signup", signUpHandler(db))
	r.HandleFunc("/user", getUserHandler(db)).Methods("GET")
	r.HandleFunc("/logout", logOutHandler(db))
}