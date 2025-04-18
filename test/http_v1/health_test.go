package integration

import (
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"boilerplate/internal/application"
	httpv1 "boilerplate/internal/infrastructure/api/http_v1"
	"boilerplate/internal/infrastructure/persistence/postgres"
	"boilerplate/internal/infrastructure/system"

	"github.com/stretchr/testify/require"
)

func (s *httpV1Suite) TestHealthEndpoint() {
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

	t.Run("should return 200 status code", func(t *testing.T) {
		resp, err := http.Get(server.URL + "/health")
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode)
	})
}
