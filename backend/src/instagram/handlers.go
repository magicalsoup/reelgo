package instagram

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"database/sql"

	"github.com/magicalsoup/reelgo/src/gcs"
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

func webhookHandler(db *sql.DB) http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			fmt.Println("recieved body")

			body, err := io.ReadAll(r.Body)

			if err != nil {
				fmt.Println(err)
			}

			defer r.Body.Close()

			verify_payload := verifyReqSignature(r, body)

			if verify_payload != nil {
				http.Error(w, verify_payload.Error(), http.StatusBadRequest)
				return
			}

			var reqBody MessageWebhookObject

			if err := json.NewDecoder(bytes.NewReader(body)).Decode(&reqBody); err != nil {
				http.Error(w, "could not parse json\n" + err.Error(), http.StatusBadRequest)
				return
			}

			reel_url := reqBody.Entry[0].Messaging[0].Message.Attachments[0].Payload.Url

			fmt.Println(reel_url)

			if reel_url == "" { // emptry url
				http.Error(w, "could not get a reel_url from message" + err.Error(), http.StatusBadRequest)
				return
			}

			attraction, err := gcs.TransformVideoData(reel_url)
			if err != nil {
				fmt.Println(err.Error())
			}

			fmt.Println(attraction.Name + " " + attraction.Location)
			w.WriteHeader(http.StatusOK)

			// TODO add the attraction to the database

		} else if r.Method == http.MethodGet {
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
}
