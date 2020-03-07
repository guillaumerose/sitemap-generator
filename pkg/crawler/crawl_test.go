package crawler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestCrawlWebsite(t *testing.T) {
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.HTML(http.StatusOK, `<a href="/depth1">depth1</a>  <a href="/about">about</a> <a href="/">self</a>`)
	})
	e.GET("/about", func(c echo.Context) error {
		return c.HTML(http.StatusOK, `content`)
	})
	e.GET("/depth1", func(c echo.Context) error {
		return c.HTML(http.StatusOK, `<a href="/depth1/depth2">depth2</a>`)
	})
	e.GET("/depth1/depth2", func(c echo.Context) error {
		return c.HTML(http.StatusOK, `content`)
	})

	target := httptest.NewServer(e)
	defer target.Close()

	crawler := New(2)
	crawler.Crawl(target.URL, 10)
	crawler.Wait()
	links := crawler.VisitedURLs()
	assert.Equal(t, links, []string{"/", "/about", "/depth1", "/depth1/depth2"})
}

func TestDiscardErrorPages(t *testing.T) {
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.HTML(http.StatusOK, `<a href="/about">about</a>`)
	})
	e.GET("/about", func(c echo.Context) error {
		return c.HTML(http.StatusNotFound, "Not Found!")
	})

	target := httptest.NewServer(e)
	defer target.Close()

	crawler := New(1)
	crawler.Crawl(target.URL, 10)
	crawler.Wait()
	links := crawler.VisitedURLs()
	assert.Equal(t, links, []string{"/"})
}
