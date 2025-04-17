package integration

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/docker/go-connections/nat"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type httpV1Suite struct {
	suite.Suite
	httpServer testcontainers.Container
	config     httpV1SuiteConfig
}

type httpV1SuiteConfig struct {
	HTTPPort       string
	requestTimeout time.Duration
	serverURL      string
}

func TestHttpV1Suite(t *testing.T) {
	suite.Run(t, new(httpV1Suite))
}

func (s *httpV1Suite) SetupSuite() {
	t := s.T()

	err := s.buildHTTPServer(t)
	assert.NoError(t, err)

}

func (s *httpV1Suite) buildHTTPServer(t *testing.T) error {
	ctx := context.Background()

	config := httpV1SuiteConfig{
		HTTPPort:       "8080",
		requestTimeout: 10 * time.Second,
	}

	req := testcontainers.ContainerRequest{
		FromDockerfile: testcontainers.FromDockerfile{
			Context:    "../..",
			Dockerfile: "Dockerfile",
		},
		// Use host networking instead of port mapping
		NetworkMode: "host",
		// No need for ExposedPorts when using host networking
		WaitingFor: wait.ForHTTP("/health").WithPort(nat.Port(config.HTTPPort)),
		Name:       "http-v1",
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	assert.NoError(t, err)

	s.httpServer = container

	host, err := container.Host(ctx)
	assert.NoError(t, err)

	port, err := container.MappedPort(ctx, nat.Port(config.HTTPPort))
	assert.NoError(t, err)

	s.config.serverURL = fmt.Sprintf("http://%s:%s", host, port.Port())

	return nil
}
