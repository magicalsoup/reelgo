package users

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
)

// sign up
// user supplies id and pw
// pw gets hashed and send to server
// pw gets salted with a unique salt, then hashed again and stored in db
// returns a bearer token


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

		json_resp, err := json.Marshal(token)

		if err != nil {
			http.Error(w, "could not serialize body\n" + err.Error(), http.StatusInternalServerError)
			return	
		}
		//w.WriteHeader(http.StatusCreated)
		fmt.Println(w.Header())
		w.Header().Set("Content-Type", "application/json")
		w.Write(json_resp)
		//w.Write([]byte("what is happening"))
	}
}