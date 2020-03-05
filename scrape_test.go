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
		return c.String(http.StatusOK, `... html ... <a href="/about">about</a> <a href="/">self</a> ... html ...`)
	})

	target := httptest.NewServer(e)
	defer target.Close()

	links, err := scrape(target.URL)
	require.NoError(t, err)
	assert.Equal(t, links, []string{"/about", "/"})
}

func TestScrapeBrokenPage(t *testing.T) {
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusForbidden, ``)
	})

	target := httptest.NewServer(e)
	defer target.Close()

	_, err := scrape(target.URL)
	assert.EqualError(t, err, "unexpected status code, got 403")
}

func TestCrawlWebsite(t *testing.T) {
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, `<a href="/depth1">depth1</a>  <a href="/about">about</a> <a href="/">self</a>`)
	})
	e.GET("/about", func(c echo.Context) error {
		return c.String(http.StatusOK, `content`)
	})
	e.GET("/depth1", func(c echo.Context) error {
		return c.String(http.StatusOK, `<a href="/depth1/depth2">depth2</a>`)
	})
	e.GET("/depth1/depth2", func(c echo.Context) error {
		return c.String(http.StatusOK, `content`)
	})

	target := httptest.NewServer(e)
	defer target.Close()

	links, err := crawl(target.URL, 10)
	require.NoError(t, err)
	assert.Equal(t, links, []string{"/", "/about", "/depth1", "/depth1/depth2"})
}
