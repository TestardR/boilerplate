package integration

import (
	"boilerplate/internal/domain/shared"
	"boilerplate/internal/domain/user/model"
	"boilerplate/internal/infrastructure/persistence/postgres"
	testshared "boilerplate/test-shared"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func (s *postgresSuite) TestUserStore() {
	t := s.T()

	ctx := t.Context()

	fixedClock := testshared.NewFixedClock(time.Now())
	userStore := postgres.NewUserStore(s.db)

	now := shared.OccurredAtFrom(fixedClock.Now())

	userID := model.NewID(uuid.New())
	username := "test-user"

	user := model.NewUser(userID, username, now.AsTime())
	err := userStore.Persist(ctx, user)
	assert.NoError(t, err)

	loadedUser, err := userStore.Load(ctx, userID)
	assert.NoError(t, err)
	assert.Equal(t, user, loadedUser)
}
