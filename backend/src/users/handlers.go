package users

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/magicalsoup/reelgo/.gen/reelgo/public/model"
)

// sign up
// user supplies id and pw
// pw gets hashed and send to server
// pw gets salted with a unique salt, then hashed again and stored in db

func loginHandler(db *sql.DB) http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {
		data := &UserDataPayload{}
		err := json.NewDecoder(r.Body).Decode(data)

		defer r.Body.Close()

		if err != nil {
			http.Error(w, "could not parse request body into json\n" + err.Error(), http.StatusBadRequest)
			return 
		}

		user, err := getUser(db, data.Email)

		if err != nil {
			http.Error(w, "something went wrong\n" + err.Error(), http.StatusInternalServerError)
			return
		}

		if user == nil { // no user found, should sign up user
			// TODO probably return like a new resource response so frontend can redirect user to sign up
			w.WriteHeader(http.StatusNotFound)
			return
		}

		authenticated := loginUser(user, data.Hashed_password)

		if !authenticated { // user supplied wrong password
			w.WriteHeader(http.StatusUnauthorized)
		}

		token, err := refreshSessionToken(db, user.UID)

		if err != nil {
			http.Error(w, "something went wrong\n" + err.Error(), http.StatusInternalServerError)
			return
		}

		// otherwise write 200 
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(token)
	}
}

func setAuthCookies(w *http.ResponseWriter, token *model.Tokens) {
	id_cookie := http.Cookie{
		Name: "user-id", 
		Value: strconv.Itoa(int(*token.UID)), 
		Path: "/",
		MaxAge: int(*token.ExpiryTime),
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}

	bearer_cookie := http.Cookie{
		Name: "user-bearer-token",
		Value: *token.BearerToken,
		Path: "/",
		MaxAge: int(*token.ExpiryTime),
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}

	http.SetCookie(*w, &id_cookie)
	http.SetCookie(*w, &bearer_cookie)
}

func signUpHandler(db *sql.DB) http.HandlerFunc {
	return func (w http.ResponseWriter, r* http.Request) {

		data := &UserDataPayload{}
		err := json.NewDecoder(r.Body).Decode(data)

		defer r.Body.Close()

		if err != nil {
			http.Error(w, "could not parse request body into json\n" + err.Error(), http.StatusBadRequest)
			return 
		}

		user, err := signUpUser(db, data.Email, data.Hashed_password)

		if err != nil {
			http.Error(w, "something went wrong\n" + err.Error(), http.StatusInternalServerError)
			return
		}

		token, err := createSessionToken(db, user.UID)

		
		if err != nil {
			http.Error(w, "something went wrong\n" + err.Error(), http.StatusInternalServerError)
			return
		}

		setAuthCookies(&w, token);
		w.WriteHeader(http.StatusCreated) // write this after setting cookies always
	}
}