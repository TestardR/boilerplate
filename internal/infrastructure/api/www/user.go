package www

import (
	"time"

	user "boilerplate/internal/domain/user/model"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	UpdatedAt time.Time `json:"updated_at"`
}

func ToUser(userModel user.User) User {
	return User{
		ID:        userModel.ID(),
		Username:  userModel.Username(),
		UpdatedAt: userModel.UpdatedAt(),
	}
}
