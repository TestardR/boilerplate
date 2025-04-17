package application

import (
	"context"

	"boilerplate/internal/application/command"
	"boilerplate/internal/application/query"
	"boilerplate/internal/domain/shared"
	"boilerplate/internal/domain/user/model"
	userrepository "boilerplate/internal/domain/user/repository"
)

type UserService struct {
	userPersister userrepository.Persister
	userLoader    userrepository.Loader
	clock         shared.CurrentTime
}

func NewUserService(
	userPersister userrepository.Persister,
	userLoader userrepository.Loader,
	clock shared.CurrentTime,
) UserService {
	return UserService{
		userPersister: userPersister,
		userLoader:    userLoader,
		clock:         clock,
	}
}

func (s UserService) AddUser(ctx context.Context, cmd command.AddUser) error {
	user := model.NewUser(
		cmd.ID(),
		cmd.Username(),
		s.clock.Now(),
	)

	return s.userPersister.Persist(ctx, user)
}

func (s UserService) GetUser(ctx context.Context, qry query.GetUser) (model.User, error) {
	return s.userLoader.Load(ctx, qry.ID())
}
