package client

import (
	"fmt"
	"net/http"
	"time"
)

type Client struct {
	url  string
	http *http.Client
}

func New(url string) *Client {
	return &Client{
		url: url,
		http: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func (c *Client) Healthcheck() error {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/healthz", c.url), nil)
	if err != nil {
		return err
	}
	res, err := c.http.Do(req)
	if err != nil {
		return err
	}
	if res.StatusCode != 200 {
		return fmt.Errorf("unexpected status code, got %d", res.StatusCode)
	}
	return nil
}
