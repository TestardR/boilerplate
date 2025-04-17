package httpv1

import (
	"context"

	"boilerplate/internal/application/command"
	"boilerplate/internal/application/query"
	user "boilerplate/internal/domain/user/model"
)

//go:generate go tool go.uber.org/mock/mockgen -source=handler.go -destination=./mock/user_service.go -package=mock

type UserService interface {
	AddUser(ctx context.Context, cmd command.AddUser) error
	GetUser(ctx context.Context, qry query.GetUser) (user.User, error)
}

type Handler struct {
	userService UserService
}

func NewHandler(userService UserService) Handler {
	return Handler{
		userService: userService,
	}
}
