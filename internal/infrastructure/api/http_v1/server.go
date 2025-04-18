package httpv1

import (
	"net/http"

	"boilerplate/internal/domain/shared"
)

func NewHttServer(
	config Config,
	logger shared.Logger,
	handler Handler,
) *http.Server {
	mux := http.NewServeMux()

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		logger.DebugContext(r.Context(), "health-check")

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok")) //nolint:errcheck
	})

	mux.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			handler.CreateUser(w, r)
		case http.MethodGet:
			handler.GetUser(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	return &http.Server{
		Addr:         config.Address,
		ReadTimeout:  config.Timeout,
		WriteTimeout: config.Timeout,
		Handler:      mux,
	}
}
