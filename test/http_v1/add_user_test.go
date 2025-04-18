package integration

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"boilerplate/internal/application"
	httpv1 "boilerplate/internal/infrastructure/api/http_v1"
	"boilerplate/internal/infrastructure/persistence/postgres"
	"boilerplate/internal/infrastructure/system"
)

func (s *httpV1Suite) TestAddUserEndpoint() {
	t := s.T()

	logger := slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.Level(s.config.LogLevel),
	}))

	postgresClient, err := postgres.NewClient(s.config.Postgres)
	require.NoError(t, err)

	userStore := postgres.NewUserStore(postgresClient.DB())

	userService := application.NewUserService(
		userStore,
		userStore,
		system.NewClock(),
	)

	handler := httpv1.NewHandler(userService)

	httpServer := httpv1.NewHttServer(s.config.HTTP, logger, handler)

	server := httptest.NewServer(httpServer.Handler)
	defer server.Close()

	t.Run("should successfully create a user", func(t *testing.T) {
		requestBody := httpv1.AddUserRequest{
			Username: "testuser",
		}
		body, err := json.Marshal(requestBody)
		require.NoError(t, err)

		req, err := http.NewRequest(http.MethodPost, server.URL+"/users", bytes.NewReader(body))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, http.StatusCreated, resp.StatusCode)

		var response struct {
			ID       string `json:"id"`
			Username string `json:"username"`
		}
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		require.NotEmpty(t, response.ID)
		require.Equal(t, "testuser", response.Username)

		getResp, err := http.Get(server.URL + "/users?id=" + response.ID)
		require.NoError(t, err)
		defer getResp.Body.Close()

		require.Equal(t, http.StatusOK, getResp.StatusCode)

		var getResponse struct {
			ID       string `json:"id"`
			Username string `json:"username"`
		}
		err = json.NewDecoder(getResp.Body).Decode(&getResponse)
		require.NoError(t, err)

		require.Equal(t, response.ID, getResponse.ID)
		require.Equal(t, response.Username, getResponse.Username)
	})
}
