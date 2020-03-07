package server

import (
	"net"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Server struct {
	handler http.Handler
}

func New() *Server {
	e := echo.New()
	e.Use(middleware.Logger())
	e.GET("/healthz", func(c echo.Context) error {
		return c.String(http.StatusOK, "OK")
	})

	return &Server{
		handler: e,
	}
}

func (s *Server) Start(ln net.Listener) error {
	return http.Serve(ln, s.handler)
}
