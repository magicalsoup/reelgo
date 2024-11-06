package trips

import (
	"database/sql"

	"github.com/gorilla/mux"
)

func AddRoutes(r *mux.Router, db* sql.DB) {
	r.HandleFunc("/trips", getTripsHandler(db))
}