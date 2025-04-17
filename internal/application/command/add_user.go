package command

import "boilerplate/internal/domain/user/model"

type AddUser struct {
	id       model.ID
	username string
}

func NewAddUser(id model.ID, username string) AddUser {
	return AddUser{
		id:       id,
		username: username,
	}
}

func (c AddUser) ID() model.ID {
	return c.id
}

func (c AddUser) Username() string {
	return c.username
}
