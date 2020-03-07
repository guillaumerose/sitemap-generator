package server

import (
	"fmt"
	"net"
	"testing"

	"github.com/guillaumerose/sitemap-generator/pkg/client"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHealthCheck(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	defer listener.Close()
	server := New()
	go func() {
		if err := server.Start(listener); err != nil {
			logrus.Error(err)
		}
	}()

	client := client.New(fmt.Sprintf("http://%s", listener.Addr().String()))
	assert.NoError(t, client.Healthcheck())
}
