package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"net/http"
)

// flow user recieves !verify [hashed_id] on client
// user sends that text to reelgo.app
// reelgo.app recieves that text and stores a 6 digit verification code, the hashed id in backend and the igsid, sends 6 digit code back to user
// user recieves and types 6 digit code in frontend
// frontend recieves and sends back the 6 digit code and the user id (gotten from when the response when user signs up/logs in)

// the function below checks if the 6 digit code and hashed id matches in the database
// if so, stores the igsid to the user to link, and returns a statusOK
func codeCheckHandler(db *sql.DB) http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {
	
		data := &VerificationPayload{}
		err := json.NewDecoder(r.Body).Decode(data)
		
		defer r.Body.Close()

		if err != nil {
			http.Error(w, "could not parse request body into json\n" + err.Error(), http.StatusBadRequest)
			return 
		}

		h := hmac.New(sha256.New, []byte(string(data.Uid)))
		huid := hex.EncodeToString(h.Sum(nil))

		verification_code, err := getVerificationCode(db, huid)

		if err != nil || *verification_code.Code != data.Code{
			http.Error(w, "something went wrong \n" + err.Error(), http.StatusNotFound)
			return 
		}

		err = storeIGIDToUser(db, data.Uid, *verification_code.InstagramID)

		if err != nil {
			http.Error(w, "something went wrong\n" + err.Error(), http.StatusInternalServerError)
			return 
		}

		w.WriteHeader(http.StatusOK)
	}
}