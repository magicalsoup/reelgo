package instagram

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
)

type messagePayload struct {
	Recipient Recipient `json:"recipient"`
	Message struct {
		Text string `json:"text"`
	} `json:"message"`
}

func sendMessageToUser(ig_id string, message string) error {
	api_url := fmt.Sprintf("https://graph.facebook.com/v21.0/%s/messages", os.Getenv("APP_IGSID"))
	client := &http.Client{}
	fmt.Println(api_url)

	payload := &messagePayload{}
	payload.Recipient.Id = ig_id
	payload.Message.Text = message

	b, err := json.Marshal(payload)
	fmt.Println(string(b))

	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", api_url, bytes.NewBuffer(b))
	if err != nil {
		return err
	}

	bearer := "Bearer " + os.Getenv("INSTAGRAM_USER_ACCESS_TOKEN")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", bearer)
	fmt.Println(req.Header)

	resp, err := client.Do(req)

	if resp.StatusCode != http.StatusOK {
		return errors.New("did not get a statusOK, instead got: " + string(resp.StatusCode))
	}
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

