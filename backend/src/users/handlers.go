package users

import (
	"database/sql"
	"encoding/json"

	"net/http"
	"github.com/magicalsoup/reelgo/src/util"
)

// sign up
// user supplies id and pw
// pw gets hashed and send to server
// pw gets salted with a unique salt, then hashed again and stored in db
func loginHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := &UserAuthPayload{}
		err := json.NewDecoder(r.Body).Decode(data)

		defer r.Body.Close()

		if err != nil {
			http.Error(w, "could not parse request body into json\n"+err.Error(), http.StatusBadRequest)
			return
		}

		user, err := GetUserByEmail(db, data.Email)

		if err != nil {
			http.Error(w, "something went wrong\n"+err.Error(), http.StatusInternalServerError)
			return
		}

		if user == nil { // no user found, should sign up user
			// TODO probably return like a new resource response so frontend can redirect user to sign up
			w.WriteHeader(http.StatusNotFound)
			return
		}

		hashed_secret := getHashedPassword(data.Hashed_password, user.Salt)

		if hashed_secret != user.HashedPassword { // user supplied wrong password
			w.WriteHeader(http.StatusUnauthorized)
		}

		token, err := RefreshSessionToken(db, user.UID)

		if err != nil {
			http.Error(w, "something went wrong\n"+err.Error(), http.StatusInternalServerError)
			return
		}

		// otherwise write 200
		util.SetAuthCookies(w, token)
		w.WriteHeader(http.StatusOK)
	}
}

func signUpHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		data := &UserAuthPayload{}
		err := json.NewDecoder(r.Body).Decode(data)

		defer r.Body.Close()

		if err != nil {
			http.Error(w, "could not parse request body into json\n"+err.Error(), http.StatusBadRequest)
			return
		}

		user, err := CreateUser(db, data.Name, data.Email, data.Hashed_password)

		if err != nil {
			http.Error(w, "something went wrong\n"+err.Error(), http.StatusInternalServerError)
			return
		}

		token, err := CreateSessionToken(db, user.UID)

		if err != nil {
			http.Error(w, "something went wrong\n"+err.Error(), http.StatusInternalServerError)
			return
		}

		util.SetAuthCookies(w, token)
		w.WriteHeader(http.StatusCreated) // write this after setting cookies always
	}
}

func logOutHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookies := r.Cookies()

		bearer_token, err := util.GetBearerToken(cookies)

		if err != nil {
			http.Error(w, "user not signed in, error: "+err.Error(), http.StatusUnauthorized)
			return
		}

		err = InvalidateSessionToken(db, bearer_token)

		if err != nil {
			http.Error(w, "something went wrong", http.StatusInternalServerError)
			return
		}
	}
}

func getUserHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookies := r.Cookies()

		bearer_token, err := util.GetBearerToken(cookies)

		if err != nil {
			http.Error(w, "user not signed in error: "+err.Error(), http.StatusUnauthorized)
			return
		}

		user, err := GetUserByToken(db, bearer_token)

		if err != nil {
			http.Error(w, "user not found or invalid token", http.StatusUnauthorized)
			return
		}

		// kind of a hacky solution, will refine later

		ig_id := ""
		if user.InstagramID != nil {
			ig_id = *user.InstagramID
		}

		data := UserDataPayload{
			UID:          user.UID,
			Name:         user.Name,
			Email:        user.Email,
			Instagram_id: ig_id,
			Verified:     user.Verified,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(data)
	}
}
