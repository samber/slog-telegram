package slogtelegram

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"sync"
	"testing"
	"time"
)

type mockHTTPClient struct {
	mu             sync.Mutex
	requestHistory []http.Request
}

func (c *mockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	c.mu.Lock()
	c.requestHistory = append(c.requestHistory, *req)
	c.mu.Unlock()

	var body string
	switch {
	case strings.Contains(req.URL.Path, "/getMe"):
		body = `{"ok":true,"result":{"id":1,"is_bot":true,"first_name":"test","username":"test_bot"}}`
	case strings.Contains(req.URL.Path, "/sendMessage"):
		body = `{"ok":true,"result":{"message_id":1,"date":0,"chat":{"id":1}}}`
	default:
		body = `{"ok":true,"result":{}}`
	}

	return &http.Response{
		StatusCode: http.StatusOK,
		Header:     make(http.Header),
		Body:       io.NopCloser(strings.NewReader(body)),
	}, nil
}

func (c *mockHTTPClient) Requests() []http.Request {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.requestHistory[:]
}

func TestNewTelegramHandler(t *testing.T) {
	const customEndpoint = "https://proxy.example.com/mtproto"

	client := &mockHTTPClient{}

	handler := Option{
		Token:       "123456:ABC-DEF",
		Username:    "@test_channel",
		HTTPClient:  client,
		APIEndpoint: customEndpoint,
	}.NewTelegramHandler()
	require.NotNil(t, handler)

	requests := client.Requests()

	require.Len(t, requests, 1)
	assert.True(t, strings.HasPrefix(requests[0].URL.String(), customEndpoint), "GetMe should be requested through the custom API endpoint")

	logger := slog.New(handler)
	logger.Info("test message")

	// Handle() sends the message asynchronously, so wait for it
	deadline := time.Now().Add(time.Second)
	for len(client.Requests()) < 2 && time.Now().Before(deadline) {
		time.Sleep(10 * time.Millisecond)
	}

	requests = client.Requests()
	require.Len(t, requests, 2)

	assert.True(t, strings.HasPrefix(requests[1].URL.String(), customEndpoint), "sendMessage request should go go through the custom API endpoint")
	assert.True(t, strings.Contains(requests[1].URL.Path, "/sendMessage"))
}
