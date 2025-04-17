package entity

import (
	"time"

	"boilerplate/internal/domain/user/model"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `db:"user_id"`
	Username  string    `db:"username"`
	UpdatedAt time.Time `db:"updated_at"`
}

func UserFromEntity(user User) model.User {
	return model.NewUser(
		model.NewID(user.ID),
		user.Username,
		user.UpdatedAt,
	)
}
