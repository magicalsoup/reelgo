package trips

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/magicalsoup/reelgo/src/users"
	"github.com/magicalsoup/reelgo/src/util"
)

// gets all trips and their attractions
func getTripsHandler(db *sql.DB) http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {
		cookies := r.Cookies()

		bearer_token, err := util.GetBearerToken(cookies)

		if err != nil {
			http.Error(w, "user not signed in error: "+err.Error(), http.StatusUnauthorized)
			return
		}

		user, err := users.GetUserByToken(db, bearer_token)

		if err != nil {
			http.Error(w, "user not found or invalid token", http.StatusUnauthorized)
			return
		}

		trips, err := GetTrips(db, user)

		if err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusAccepted)
		json.NewEncoder(w).Encode(trips)
	}
}
