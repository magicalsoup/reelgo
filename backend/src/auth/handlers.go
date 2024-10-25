package auth

import (
	"database/sql"
	"net/http"
)

// flow user recieves !verify [hashed_id] on client
// user sends that text to reelgo.app
// reelgo.app recieves that text and stores a 6 digit verification code, the hashed id in backend and the igsid, sends 6 digit code back to user
// user recieves and types 6 digit code in frontend
// frontend recieves and sends back the 6 digit code and the user id (gotten from when the response when user signs up/logs in)
// the function below checks if the 6 digit code and hashed id matches in the database
// if so, stores the igsid to the user to link, and returns a statusOK

func verificationHandler(db *sql.DB) http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {
		
	}
}