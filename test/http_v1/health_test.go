package integration

import (
	"net/http"

	"github.com/stretchr/testify/assert"
)

func (s *httpV1Suite) TestHealthEndpoint() {
	t := s.T()

	client := &http.Client{
		Timeout: s.config.requestTimeout,
	}

	resp, err := client.Get(s.config.serverURL + "/health")
	assert.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}
