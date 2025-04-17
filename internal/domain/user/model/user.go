package model

import (
	"time"
)

type User struct {
	id        ID
	username  string
	updatedAt time.Time
}

func NewUser(id ID, username string, updatedAt time.Time) User {
	return User{
		id:        id,
		username:  username,
		updatedAt: updatedAt,
	}
}

func (u User) ID() ID {
	return u.id
}

func (u User) Username() string {
	return u.username
}

func (u User) UpdatedAt() time.Time {
	return u.updatedAt
}
