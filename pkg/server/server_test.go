package server

import (
	"fmt"
	"net"
	"testing"

	"github.com/guillaumerose/sitemap-generator/pkg/client"
	"github.com/guillaumerose/sitemap-generator/pkg/types"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHealthCheck(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	defer listener.Close()
	server := New(&mockRepository{})
	go func() {
		if err := server.Start(listener); err != nil {
			logrus.Error(err)
		}
	}()

	client := client.New(fmt.Sprintf("http://%s", listener.Addr().String()))
	assert.NoError(t, client.Healthcheck())
}

func TestCreateGet(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	defer listener.Close()
	server := New(&mockRepository{})
	go func() {
		if err := server.Start(listener); err != nil {
			logrus.Error(err)
		}
	}()

	client := client.New(fmt.Sprintf("http://%s", listener.Addr().String()))

	crawl, err := client.CreateCrawl(&types.Crawl{
		Spec: types.CrawlSpec{
			URL: "https://www.redhat.com",
		},
	})
	assert.NoError(t, err)
	assert.Equal(t, crawl.ID, "1")
	assert.Equal(t, crawl.Status.Done, false)

	crawl, err = client.GetCrawl("1")
	assert.NoError(t, err)
	assert.Equal(t, crawl.Status.Done, true)

	links, err := client.GetCrawlLinks("1")
	assert.NoError(t, err)
	assert.Equal(t, links, []string{"link1", "link2"})
}

type mockRepository struct{}

func (mockRepository) Create(req *types.Crawl) (*types.Crawl, error) {
	return &types.Crawl{
		ID: "1",
	}, nil
}

func (mockRepository) Get(id string) (*types.Crawl, error) {
	return &types.Crawl{
		ID: "1",
		Status: types.CrawlStatus{
			Done: true,
			Size: 42,
		},
	}, nil
}

func (mockRepository) GetLinks(id string) ([]string, error) {
	return []string{"link1", "link2"}, nil
}
