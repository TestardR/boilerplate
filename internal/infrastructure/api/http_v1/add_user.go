package httpv1

import (
	"encoding/json"
	"net/http"

	"boilerplate/internal/application/command"
)

type AddUserRequest struct {
	Username string `json:"username" validate:"required,min=3"`
}

func (h Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	if contentType := r.Header.Get("Content-Type"); contentType != "application/json" {
		http.Error(w, "Content-Type must be application/json", http.StatusUnsupportedMediaType)
		return
	}

	var request AddUserRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if request.Username == "" {
		http.Error(w, "Username is required", http.StatusBadRequest)
		return
	}
	if len(request.Username) < 3 {
		http.Error(w, "Username must be at least 3 characters", http.StatusBadRequest)
		return
	}

	cmd := command.NewAddUser(request.Username)
	err := h.userService.AddUser(r.Context(), cmd)
	if err != nil {
		http.Error(w, "Failed to add user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
