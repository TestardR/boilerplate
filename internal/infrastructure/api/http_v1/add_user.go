package httpv1

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/google/uuid"

	"boilerplate/internal/application/command"
	"boilerplate/internal/application/query"
	usererrors "boilerplate/internal/domain/user"
	"boilerplate/internal/domain/user/model"
	"boilerplate/internal/infrastructure/api/www"
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

	cmd := command.NewAddUser(model.NewID(uuid.New()), request.Username)
	err := h.userService.AddUser(r.Context(), cmd)
	if err != nil {
		http.Error(w, "Failed to add user", http.StatusInternalServerError)
		return
	}

	user, err := h.userService.GetUser(r.Context(), query.NewGetUser(cmd.ID()))
	if err != nil && errors.Is(err, usererrors.ErrUserNotFound) {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "Failed to get user", http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(www.ToUser(user))
	if err != nil {
		http.Error(w, "Failed to marshal user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	_, _ = w.Write(data) //nolint:errcheck
}
