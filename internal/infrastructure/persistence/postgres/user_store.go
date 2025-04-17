package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	usererrors "boilerplate/internal/domain/user"
	user "boilerplate/internal/domain/user/model"
	"boilerplate/internal/infrastructure/persistence/postgres/entity"

	"github.com/jmoiron/sqlx"
)

type UserStore struct {
	db *sqlx.DB
}

func NewUserStore(db *sqlx.DB) UserStore {
	return UserStore{db: db}
}

func (s UserStore) Load(ctx context.Context, id user.ID) (user.User, error) {
	var userEntity entity.User

	err := s.db.GetContext(ctx, &userEntity,
		`SELECT 
			id, 
			username, 
			updated_at 
		FROM users 
		WHERE user_id = $1`,
		id.ID(),
	)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return user.User{}, usererrors.ErrUserNotFound
	}

	if err != nil {
		return user.User{}, fmt.Errorf("failed to load user: %w", err)
	}

	return entity.UserFromEntity(userEntity), nil
}

func (s UserStore) Persist(ctx context.Context, userModel user.User) error {
	result, err := s.db.ExecContext(ctx,
		`INSERT INTO users (
			id, 
			username, 
			updated_at
		) VALUES (
			$1, 
			$2, 
			$3
		)`,
		userModel.ID(),
		userModel.Username(),
		userModel.UpdatedAt(),
	)
	if err != nil {
		return fmt.Errorf("failed to persist user: %w", err)
	}

	affectedRows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}
	if affectedRows == 0 {
		return fmt.Errorf("expected user id %q to be inserted", userModel.ID().ID())
	}

	return nil
}
