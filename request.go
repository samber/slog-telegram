package slogtelegram

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// TODO: does a request to see if token is correct
//
// https://api.telegram.org/bot<token>/getMe
func (o *Option) checkInit() error {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/getMe", o.Token)

	return o.makeRequest("GET", url, "")
}

// https://core.telegram.org/bots/api#available-types
//
//	curl -X POST \
//	     -H 'Content-Type: application/json' \
//	     -d '{"chat_id": "<your-chat-id>", "text": "This is a test from curl", "disable_notification": true}' \
//	     https://api.telegram.org/bot<your-bot-token>/sendMessage
func (o *Option) sendMessage(msg string) error {

	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", o.Token)

	value := struct {
		ChatId              string `json:"chat_id"`
		Text                string `json:"text"`
		ParseMode           string `json:"parse_mode,omitempty"`
		DisableNotification bool   `json:"disable_notification"`
	}{
		ChatId:              o.ChatId,
		Text:                msg,
		ParseMode:           o.ParseMode,
		DisableNotification: false, // TODO: make it configurable
	}

	bytes, err := json.Marshal(value)

	if err != nil {
		return fmt.Errorf("slog-telegram: failed to send log: %s", o.redactSensitiveInfo(err.Error()))
	}

	return o.makeRequest("POST", url, string(bytes))
}

func (o *Option) makeRequest(method, url, body string) (err error) {

	// var reader *strings.Reader

	// if body == "" {
	// 	reader = nil // have to be explicit
	// } else {
	// 	reader = strings.NewReader(body)

	// }

	// fmt.Printf("making request %s %s %s\n", method, url, body)

	// req, err := http.NewRequest(method, url, reader)
	// if err != nil {
	// 	return fmt.Errorf("slog-telegram: failed to send log: %s", o.redactSensitiveInfo(err.Error()))
	// }

	// req.Header.Set("Content-Type", "application/json")
	// resp, err := o.HttpClient.Do(req)
	// if err != nil {
	// 	return fmt.Errorf("slog-telegram: failed to send log: %s", o.redactSensitiveInfo(err.Error()))
	// }
	var resp *http.Response
	if method == "GET" {
		resp, err = http.Get(url)
		if err != nil {
			return err
		}
	} else if method == "POST" {
		resp, err = http.Post(url, "application/json", strings.NewReader(body))
		if err != nil {
			return err
		}
	} else {
		panic("not implemented")
	}

	if resp.StatusCode != http.StatusOK {
		errBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		return fmt.Errorf("slog-telegram: failed to send log: [%s] %s", resp.Status, string(errBytes))
	}

	return nil
}

// this does not handle error wrapping
func (o *Option) redactSensitiveInfo(errMsg string) string {
	errMsg = strings.ReplaceAll(errMsg, o.Token, "<TOKEN>")
	errMsg = strings.ReplaceAll(errMsg, o.ChatId, "<CHATID>")

	return errMsg
}
