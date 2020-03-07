package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestScrapeOnePage(t *testing.T) {
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.HTML(http.StatusOK, `... html ... <a href="/about">about</a> <a href="/">self</a> <a href="//external">external</a> <a href="http://external">external</a> ... html ...`)
	})

	target := httptest.NewServer(e)
	defer target.Close()

	scraper := scraper{
		client: &http.Client{},
	}
	links, err := scraper.scrape(target.URL)
	require.NoError(t, err)
	assert.Equal(t, links, []string{"/about", "/"})
}

func TestScrapeBrokenPage(t *testing.T) {
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.HTML(http.StatusForbidden, ``)
	})

	target := httptest.NewServer(e)
	defer target.Close()

	scraper := scraper{
		client: &http.Client{},
	}
	_, err := scraper.scrape(target.URL)
	assert.EqualError(t, err, "unexpected status code, got 403")
}

func TestCleanURLs(t *testing.T) {
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.HTML(http.StatusOK, `... html ... <a href="/page1?foo=bar">page1</a> <a href="/page2#anchor">page2</a> <a href="/page3/">page3</a> ... html ...`)
	})

	target := httptest.NewServer(e)
	defer target.Close()

	scraper := scraper{
		client: &http.Client{},
	}
	links, err := scraper.scrape(target.URL)
	require.NoError(t, err)
	assert.Equal(t, links, []string{"/page1", "/page2", "/page3"})
}
