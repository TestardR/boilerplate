package query

import (
	user "boilerplate/internal/domain/user/model"
)

type GetUser struct {
	id user.ID
}

func NewGetUser(id user.ID) GetUser {
	return GetUser{
		id: id,
	}
}

func (q GetUser) ID() user.ID {
	return q.id
}
