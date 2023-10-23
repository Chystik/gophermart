package httpclient

import (
	"net/http"
	"time"
)

const (
	defaultTimeout = 20 * time.Second
)

type Client struct {
	client  *http.Client
	timeout time.Time
}

func NewClient(opts ...Options) *Client {
	httpClietn := &http.Client{
		Timeout: defaultTimeout,
	}

	client := &Client{client: httpClietn}

	for _, opt := range opts {
		opt(client)
	}

	return client
}
