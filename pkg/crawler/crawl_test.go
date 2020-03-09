package crawler

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/guillaumerose/sitemap-generator/pkg/types"
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

	crawler := New(types.CrawlSpec{
		URL:         target.URL,
		MaxDepth:    10,
		Parallelism: 2,
	})
	crawler.Crawl()
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

	crawler := New(types.CrawlSpec{
		URL:         target.URL,
		MaxDepth:    10,
		Parallelism: 2,
	})
	crawler.Crawl()
	crawler.Wait()
	links := crawler.VisitedURLs()
	assert.Equal(t, links, []string{"/"})
}

func TestDepthLimiter(t *testing.T) {
	target := httptest.NewServer(deepWebsite())
	defer target.Close()

	crawler := New(types.CrawlSpec{
		URL:         target.URL,
		MaxDepth:    5,
		Parallelism: 2,
	})
	crawler.Crawl()
	crawler.Wait()
	links := crawler.VisitedURLs()
	assert.Equal(t, links, []string{"/", "/2", "/3", "/4", "/5"})
}

func TestNoDepthLimit(t *testing.T) {
	target := httptest.NewServer(deepWebsite())
	defer target.Close()

	crawler := New(types.CrawlSpec{
		URL:         target.URL,
		MaxDepth:    -1,
		Parallelism: 2,
	})
	crawler.Crawl()
	crawler.Wait()
	links := crawler.VisitedURLs()
	assert.Equal(t, links, []string{"/", "/10", "/2", "/3", "/4", "/5", "/6", "/7", "/8", "/9"})
}

func deepWebsite() *echo.Echo {
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.HTML(http.StatusOK, `<a href="/2">2</a>`)
	})
	e.GET("/:depth", func(c echo.Context) error {
		depth, err := strconv.Atoi(c.Param("depth"))
		if err != nil {
			return err
		}
		next := min(10, depth+1)
		return c.HTML(http.StatusOK, fmt.Sprintf(`<a href="/%d">%d</a>`, next, next))
	})
	return e
}

func min(a, b int) int {
	if a > b {
		return b
	}
	return a
}
