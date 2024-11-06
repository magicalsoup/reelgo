package instagram

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/magicalsoup/reelgo/src/auth"
	"github.com/magicalsoup/reelgo/src/gcs"
	"github.com/magicalsoup/reelgo/src/trips"
	"github.com/magicalsoup/reelgo/src/users"
)

func verifyReqSignature(r *http.Request, buffer []byte) error {
	signature := r.Header.Get("x-hub-signature-256")
	if signature == "" {
		fmt.Println("Couldn't find x-hub-signature-256 in headers.")
	} else {
		elements := strings.Split(signature, "=")
		signature_hash := elements[1]

		h := hmac.New(sha256.New, []byte(os.Getenv("APP_SECRET")))
		h.Write(buffer)

		expected_hash := hex.EncodeToString(h.Sum(nil))

		if expected_hash != signature_hash {
			return errors.New("couldn't validate the request signature")
		}
	}
	return nil
}

func verifyWebHookHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		mode := r.URL.Query().Get("hub.mode")
		token := r.URL.Query().Get("hub.verify_token")
		challenge := r.URL.Query().Get("hub.challenge")

		if mode != "" && token != "" {
			if mode == "subscribe" && token == os.Getenv("VERIFY_TOKEN") {
				fmt.Println("WEBHOOK_VERIFIED")
				w.WriteHeader(http.StatusOK) // send back the 200 status back to request
				w.Write([]byte(challenge))   // sends back the challenge token
			} else {
				w.WriteHeader(http.StatusForbidden) // send back the 403 status back to request
			}
		}
	}
}

func messageWebhookHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		
		body, err := io.ReadAll(r.Body)
		defer r.Body.Close()

		if err != nil {
			http.Error(w, "could not read request body\n", http.StatusBadRequest)
			return
		}
		
		verify_payload := verifyReqSignature(r, body)

		if verify_payload != nil {
			http.Error(w, verify_payload.Error(), http.StatusBadRequest)
			return
		}

		var reqBody MessageWebhookObject

		if err := json.NewDecoder(bytes.NewReader(body)).Decode(&reqBody); err != nil {
			http.Error(w, "could not parse json\n", http.StatusBadRequest)
			return
		}

		text := reqBody.Entry[0].Messaging[0].Message.Text
		ig_id := reqBody.Entry[0].Messaging[0].Sender.Id
		recipient_id := reqBody.Entry[0].Messaging[0].Recipient.Id

		fmt.Println("sender id", ig_id)
		fmt.Println("recipient id", recipient_id)
		fmt.Println(text)

		// this means that the message was not sent to this app
		if recipient_id != os.Getenv("APP_ID") {
			return
		}

		user, err := users.GetUserByInstagramID(db, ig_id)
		
		if err != nil {
			http.Error(w, "something went wrong\n" + err.Error(), http.StatusInternalServerError)
			return
		}

		text_parts := strings.Split(text, ":")

		// expects ![VERIFY_COMMAND]:[user_id]
		// user wants to verify their account (and not already verified)
		if len(text_parts) == 2 && strings.ToLower(text_parts[0]) == os.Getenv("VERIFY_COMMAND") && !user.Verified { 
			
			uidStr := text_parts[1] // hashed id is the rest of the code
			authcode := auth.Generate6DigitCode()

			uid, converr := strconv.Atoi(uidStr)

			if converr != nil {
				fmt.Println("could not parse message\n" + converr.Error())
				return
			}

			dberr := addVerificationCodeToDB(db, authcode, int32(uid), ig_id)
			if dberr != nil {
				fmt.Println("could not add verification code to db\n" + dberr.Error())
				http.Error(w, "could not add verification code to db\n" + dberr.Error(), http.StatusBadRequest)
				return
			}

			err := sendMessageToUser(ig_id, authcode)

			if err != nil {
				fmt.Println("could not send user message\n"+err.Error())
				http.Error(w, "could not send user message\n"+err.Error(), http.StatusBadRequest)
				return
			}
			return
		}

		if !user.Verified { // not verified/linked user, can't save reels
			return 
		}

		if len(reqBody.Entry[0].Messaging[0].Message.Attachments) == 0 {
			return // no attachments
		}

		reel_url := reqBody.Entry[0].Messaging[0].Message.Attachments[0].Payload.Url

		if reel_url == "" { // emptry url
			http.Error(w, "could not get a reel_url from message", http.StatusBadRequest)
			return
		}

		attraction, err := gcs.TransformVideoData(reel_url)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		err = trips.AddAttraction(db, attraction, user)

		if err != nil {
			fmt.Println(err.Error())
			return
		}
	
		w.WriteHeader(http.StatusOK)
		
	}
}
