package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/guillaumerose/sitemap-generator/pkg/types"
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
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code, got %d", res.StatusCode)
	}
	return nil
}

func (c *Client) CreateCrawl(cr *types.Crawl) (*types.Crawl, error) {
	bin, err := json.Marshal(cr)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/crawls", c.url), bytes.NewReader(bin))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	res, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("unexpected status code, got %d", res.StatusCode)
	}
	defer res.Body.Close()
	bin, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	var crawl types.Crawl
	if err := json.Unmarshal(bin, &crawl); err != nil {
		return nil, err
	}
	return &crawl, nil
}

func (c *Client) GetCrawl(id string) (*types.Crawl, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/crawls/%s", c.url, id), nil)
	if err != nil {
		return nil, err
	}
	res, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code, got %d", res.StatusCode)
	}
	defer res.Body.Close()
	bin, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	var crawl types.Crawl
	if err := json.Unmarshal(bin, &crawl); err != nil {
		return nil, err
	}
	return &crawl, nil
}

func (c *Client) GetCrawlLinks(id string) ([]string, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/crawls/%s/links", c.url, id), nil)
	if err != nil {
		return nil, err
	}
	res, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code, got %d", res.StatusCode)
	}
	defer res.Body.Close()
	bin, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	var links []string
	if err := json.Unmarshal(bin, &links); err != nil {
		return nil, err
	}
	return links, nil
}
