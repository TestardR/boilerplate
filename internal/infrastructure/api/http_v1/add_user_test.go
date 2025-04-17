package httpv1

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"boilerplate/internal/infrastructure/api/http_v1/mock"
)

func TestHandlerCreateUser(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	username := "testuser"

	tests := []struct {
		name           string
		requestBody    AddUserRequest
		contentType    string
		expect         func(*mock.MockUserService)
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "invalid content type",
			requestBody: AddUserRequest{
				Username: username,
			},
			contentType:    "text/plain",
			expect:         func(mockUserService *mock.MockUserService) {},
			expectedStatus: http.StatusUnsupportedMediaType,
			expectedBody:   "Content-Type must be application/json\n",
		},
		{
			name: "empty username",
			requestBody: AddUserRequest{
				Username: "",
			},
			contentType:    "application/json",
			expect:         func(mockUserService *mock.MockUserService) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Username is required\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			body, err := json.Marshal(tt.requestBody)
			assert.NoError(t, err)

			request := httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader(body))
			request.Header.Set("Content-Type", tt.contentType)

			recorder := httptest.NewRecorder()

			mockService := mock.NewMockUserService(ctrl)
			tt.expect(mockService)

			handler := Handler{
				userService: mockService,
			}
			handler.CreateUser(recorder, request)

			assert.Equal(t, tt.expectedStatus, recorder.Code)
			assert.Equal(t, tt.expectedBody, recorder.Body.String())
		})
	}
}
