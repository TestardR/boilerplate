package application

import (
	"context"
	"errors"
	"testing"
	"time"

	"boilerplate/internal/application/command"
	"boilerplate/internal/application/mock"
	"boilerplate/internal/application/query"
	"boilerplate/internal/domain/user/model"
	testshared "boilerplate/test-shared"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestUserServiceAddUser(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	now := time.Now()
	clock := testshared.NewFixedClock(now)

	userID := model.NewID(uuid.New())
	username := "test-user"

	tests := []struct {
		name        string
		expect      func(*mock.MockPersister)
		expectedErr error
	}{
		{
			name: "successfully add user",

			expect: func(mockPersister *mock.MockPersister) {
				expectedUser := model.NewUser(userID, username, now)
				mockPersister.EXPECT().
					Persist(ctx, expectedUser).
					Return(nil)
			},
			expectedErr: nil,
		},
		{
			name: "repository error",

			expect: func(mockPersister *mock.MockPersister) {
				expectedUser := model.NewUser(userID, username, now)
				mockPersister.EXPECT().
					Persist(ctx, expectedUser).
					Return(errors.New("repository error"))
			},
			expectedErr: errors.New("repository error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockPersister := mock.NewMockPersister(ctrl)
			tt.expect(mockPersister)

			service := UserService{
				userPersister: mockPersister,
				clock:         clock,
			}

			cmd := command.NewAddUser(userID, username)
			err := service.AddUser(ctx, cmd)

			if (tt.expectedErr != nil) || (err != nil) {
				assert.Equal(t, tt.expectedErr, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestUserServiceGetUser(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	now := time.Now()

	clock := testshared.NewFixedClock(now)

	userID := model.NewID(uuid.New())
	expectedUser := model.NewUser(userID, "testuser", now)

	tests := []struct {
		name          string
		userID        model.ID
		mockSetup     func(*mock.MockLoader)
		expectedUser  model.User
		expectedError error
	}{
		{
			name:   "successfully get user",
			userID: userID,
			mockSetup: func(mockLoader *mock.MockLoader) {
				mockLoader.EXPECT().
					Load(ctx, userID).
					Return(expectedUser, nil)
			},
			expectedUser:  expectedUser,
			expectedError: nil,
		},
		{
			name:   "user not found",
			userID: userID,
			mockSetup: func(mockLoader *mock.MockLoader) {
				mockLoader.EXPECT().
					Load(ctx, userID).
					Return(model.User{}, errors.New("user not found"))
			},
			expectedUser:  model.User{},
			expectedError: errors.New("user not found"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockPersister := mock.NewMockPersister(ctrl)
			mockLoader := mock.NewMockLoader(ctrl)
			tt.mockSetup(mockLoader)

			service := NewUserService(mockPersister, mockLoader, clock)
			qry := query.NewGetUser(tt.userID)

			user, err := service.GetUser(ctx, qry)

			if (tt.expectedError != nil) || (err != nil) {
				assert.Equal(t, tt.expectedError, err)
			} else {
				assert.Equal(t, tt.expectedUser, user)
			}
		})
	}
}
