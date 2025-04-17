package httpv1

import (
	"encoding/json"
	"errors"
	"net/http"

	"boilerplate/internal/application/query"
	usererrors "boilerplate/internal/domain/user"
	user "boilerplate/internal/domain/user/model"
	"boilerplate/internal/infrastructure/api/www"

	"github.com/google/uuid"
)

func (h Handler) GetUser(w http.ResponseWriter, r *http.Request) {
	if contentType := r.Header.Get("Content-Type"); contentType != "application/json" {
		http.Error(w, "Content-Type must be application/json", http.StatusUnsupportedMediaType)
		return
	}

	userID, err := uuid.Parse(r.URL.Query().Get("id"))
	if err != nil {
		http.Error(w, "Invalid id", http.StatusBadRequest)
		return
	}

	qry := query.NewGetUser(user.NewID(userID))
	userModel, err := h.userService.GetUser(r.Context(), qry)
	if err != nil && errors.Is(err, usererrors.ErrUserNotFound) {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "Failed to get user", http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(www.ToUser(userModel))
	if err != nil {
		http.Error(w, "Failed to marshal user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(data) // nolint:errcheck
}
