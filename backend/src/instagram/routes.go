package instagram

import (
	"database/sql"

	"github.com/gorilla/mux"
)

func AddRoutes(r *mux.Router, db* sql.DB) {
	r.HandleFunc("/webhooks", verifyWebHookHandler()).Methods("GET")
	r.HandleFunc("/webhooks", messageWebhookHandler(db)).Methods("POST")
}