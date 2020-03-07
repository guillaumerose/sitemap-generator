package main

import (
	"flag"
	"net"

	"github.com/guillaumerose/sitemap-generator/pkg/repository"
	"github.com/guillaumerose/sitemap-generator/pkg/server"
	"github.com/sirupsen/logrus"
)

var addr string

func main() {
	flag.StringVar(&addr, "addr", ":8080", "port to listen")
	flag.Parse()

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		logrus.Fatal(err)
	}
	logrus.Infof("Listening on %s", addr)
	server := server.New(repository.NewInMemoryRepository())
	if err := server.Start(listener); err != nil {
		logrus.Fatal(err)
	}
}
