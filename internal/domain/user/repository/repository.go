package repository

import (
	"context"

	user "boilerplate/internal/domain/user/model"
)

type Persister interface {
	Persist(ctx context.Context, userModel user.User) error
}

type Loader interface {
	Load(ctx context.Context, userID user.ID) (user.User, error)
}

type Updater interface {
	Update(ctx context.Context, userModel user.User) error
}
