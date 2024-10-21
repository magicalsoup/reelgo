package instagram

import (
	"github.com/gorilla/mux"
)

func AddRoutes(r *mux.Router) {
	r.HandleFunc("/webhooks", webhookHandler())
}