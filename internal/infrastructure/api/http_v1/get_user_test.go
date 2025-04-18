package httpv1

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"boilerplate/internal/application/query"
	"boilerplate/internal/domain/shared"
	usererrors "boilerplate/internal/domain/user"
	"boilerplate/internal/domain/user/model"
	"boilerplate/internal/infrastructure/api/http_v1/mock"
	"boilerplate/internal/infrastructure/api/www"
	testshared "boilerplate/test_shared"
)

func TestHandlerGetUser(t *testing.T) {
	t.Parallel()

	fixedClocked := testshared.NewFixedClock(time.Now())

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()

	now := shared.OccurredAtFrom(fixedClocked.Now())
	userID := model.NewID(uuid.New())
	username := "testuser"

	tests := []struct {
		name           string
		userIDParam    string
		contentType    string
		expect         func(*mock.MockUserService)
		expectedStatus int
		expectedBody   string
	}{
		{
			name:        "successfully get user",
			userIDParam: userID.ID().String(),
			contentType: "application/json",
			expect: func(mockUserService *mock.MockUserService) {
				expectedUser := model.NewUser(userID, username, now.AsTime())
				mockUserService.EXPECT().
					GetUser(ctx, query.NewGetUser(userID)).
					Return(expectedUser, nil)

				response := www.ToUser(expectedUser)
				_, err := json.Marshal(response)
				assert.NoError(t, err)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   fmt.Sprintf(`{"id":"%s","username":"%s"}`, userID.ID().String(), username),
		},
		{
			name:           "invalid uuid",
			userIDParam:    "invalid-uuid",
			contentType:    "application/json",
			expect:         func(mockUserService *mock.MockUserService) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Invalid id\n",
		},
		{
			name:        "user not found",
			userIDParam: userID.ID().String(),
			contentType: "application/json",
			expect: func(mockUserService *mock.MockUserService) {
				mockUserService.EXPECT().
					GetUser(ctx, query.NewGetUser(userID)).
					Return(model.User{}, usererrors.ErrUserNotFound)
			},
			expectedStatus: http.StatusNotFound,
			expectedBody:   "User not found\n",
		},
		{
			name:        "service error",
			userIDParam: userID.ID().String(),
			contentType: "application/json",
			expect: func(mockUserService *mock.MockUserService) {
				mockUserService.EXPECT().
					GetUser(ctx, query.NewGetUser(userID)).
					Return(model.User{}, errors.New("service error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "Failed to get user\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			req := httptest.NewRequest(http.MethodGet, "/users?id="+tt.userIDParam, nil)
			req.Header.Set("Content-Type", tt.contentType)

			recorder := httptest.NewRecorder()

			mockService := mock.NewMockUserService(ctrl)
			tt.expect(mockService)

			handler := Handler{
				userService: mockService,
			}
			handler.GetUser(recorder, req)

			assert.Equal(t, tt.expectedStatus, recorder.Code)
			assert.Equal(t, tt.expectedBody, recorder.Body.String())
		})
	}
}
