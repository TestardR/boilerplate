package www

import (
	user "boilerplate/internal/domain/user/model"

	"github.com/google/uuid"
)

type User struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
}

func ToUser(userModel user.User) User {
	return User{
		ID:       userModel.ID().ID(),
		Username: userModel.Username(),
	}
}
