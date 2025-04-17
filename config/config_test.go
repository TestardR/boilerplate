package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

// env test can't be parallel.
func TestLoadConfigFromEnv(t *testing.T) {
	varsToClear := getConfigEnvs(t)
	clearAllEnv(t, varsToClear)

	vars := getConfigEnvs(t)
	populateEnv(t, vars)

	config, err := FromEnv()
	require.NoError(t, err)

	require.Equal(t, 1, config.LogLevel)

	require.Equal(t, "localhost:7000", config.HTTP.Address)
	require.Equal(t, "2s", config.HTTP.Timeout.String())

	require.Equal(t, []string{"localhost:8000"}, config.EventStream.Brokers)
	require.Equal(t, "test", config.EventStream.Topic)
}

func getConfigEnvs(t *testing.T) map[string]string {
	t.Helper()

	vars := make(map[string]string)

	vars["LOG_LEVEL"] = "1"

	vars["HTTP_ADDRESS"] = "localhost:7000"
	vars["HTTP_TIMEOUT"] = "2s"

	vars["EVENT_STREAM_BROKERS"] = "localhost:8000"
	vars["EVENT_STREAM_TOPIC"] = "test"

	return vars
}

func populateEnv(t *testing.T, vars map[string]string) {
	t.Helper()
	for k, v := range vars {
		os.Setenv(k, v)
	}
}

func clearAllEnv(t *testing.T, vars map[string]string) {
	t.Helper()
	for k := range vars {
		os.Unsetenv(k)
	}
}
