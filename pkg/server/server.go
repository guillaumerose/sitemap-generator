package server

import (
	"net"
	"net/http"

	"github.com/guillaumerose/sitemap-generator/pkg/repository"
	"github.com/guillaumerose/sitemap-generator/pkg/types"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)


type Server struct {
	handler http.Handler
}

func New(repo repository.Repository) *Server {
	e := echo.New()
	e.Use(middleware.Logger())

	e.GET("/healthz", func(c echo.Context) error {
		return c.String(http.StatusOK, "OK")
	})
	e.POST("/crawls", func(c echo.Context) error {
		req := new(types.Crawl)
		if err := c.Bind(req); err != nil {
			return err
		}
		crawl, err := repo.Create(req)
		if err != nil {
			return err
		}
		return c.JSON(http.StatusCreated, crawl)
	})
	e.GET("/crawls/:id", func(c echo.Context) error {
		crawl, err := repo.Get(c.Param("id"))
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, crawl)
	})
	e.GET("/crawls/:id/links", func(c echo.Context) error {
		links, err := repo.GetLinks(c.Param("id"))
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, links)
	})

	return &Server{
		handler: e,
	}
}

func (s *Server) Start(ln net.Listener) error {
	return http.Serve(ln, s.handler)
}
