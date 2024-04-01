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

// https://core.telegram.org/method/messages.sendMessage
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
		DisableNotification bool   `json:"disable_notification"`
	}{
		ChatId:              o.ChatId,
		Text:                msg,
		DisableNotification: false, // TODO: make it configurable
	}

	bytes, err := json.Marshal(value)

	if err != nil {
		return fmt.Errorf("slog-telegram: failed to send log: %s", o.redactSensitiveInfo(err.Error()))
	}

	return o.makeRequest("POST", url, string(bytes))
}

func (o *Option) makeRequest(method, url, body string) error {
	req, err := http.NewRequest(method, url, strings.NewReader(body))
	if err != nil {
		return fmt.Errorf("slog-telegram: failed to send log: %s", o.redactSensitiveInfo(err.Error()))
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := o.HttpClient.Do(req)
	if err != nil {
		return fmt.Errorf("slog-telegram: failed to send log: %s", o.redactSensitiveInfo(err.Error()))
	}

	if resp.StatusCode != http.StatusOK {
		errBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("slog-telegram: failed to send log: %s", o.redactSensitiveInfo(err.Error()))
		}

		return fmt.Errorf("slog-telegram: failed to send log: [%s] %s", resp.Status, o.redactSensitiveInfo(string(errBytes)))
	}

	return nil
}

// this does not handle error wrapping
func (o *Option) redactSensitiveInfo(errMsg string) string {
	errMsg = strings.ReplaceAll(errMsg, o.Token, "<TOKEN>")
	errMsg = strings.ReplaceAll(errMsg, o.ChatId, "<CHATID>")

	return errMsg
}
