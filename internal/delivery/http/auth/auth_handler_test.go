package auth_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	authHandler "github.com/SergeyBogomolovv/fitflow/internal/delivery/http/auth"
	"github.com/SergeyBogomolovv/fitflow/internal/delivery/http/auth/mocks"
	"github.com/SergeyBogomolovv/fitflow/internal/domain"
	testutils "github.com/SergeyBogomolovv/fitflow/pkg/test_utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAuthHandler_Login(t *testing.T) {
	type MockBehavior func(repo *mocks.AuthService, body authHandler.LoginRequest)

	testCases := []struct {
		name           string
		body           authHandler.LoginRequest
		mockBehavior   MockBehavior
		wantStatusCode int
		wantBody       string
	}{
		{
			name: "success",
			body: authHandler.LoginRequest{
				Login:    "test_login",
				Password: "correct_password",
			},
			mockBehavior: func(repo *mocks.AuthService, body authHandler.LoginRequest) {
				repo.EXPECT().Login(mock.Anything, body.Login, body.Password).Return("valid_token", nil).Once()
			},
			wantStatusCode: http.StatusOK,
			wantBody:       `{"token":"valid_token"}` + "\n",
		},
		{
			name: "no password",
			body: authHandler.LoginRequest{
				Login: "test_login",
			},
			mockBehavior:   func(repo *mocks.AuthService, body authHandler.LoginRequest) {},
			wantStatusCode: http.StatusBadRequest,
			wantBody:       `{"status":"error","code":400,"message":"invalid payload"}` + "\n",
		},
		{
			name: "no login",
			body: authHandler.LoginRequest{
				Password: "test_password",
			},
			mockBehavior:   func(repo *mocks.AuthService, body authHandler.LoginRequest) {},
			wantStatusCode: http.StatusBadRequest,
			wantBody:       `{"status":"error","code":400,"message":"invalid payload"}` + "\n",
		},
		{
			name: "invalid credentials",
			body: authHandler.LoginRequest{
				Login:    "test_login",
				Password: "test_password",
			},
			mockBehavior: func(repo *mocks.AuthService, body authHandler.LoginRequest) {
				repo.EXPECT().Login(mock.Anything, body.Login, body.Password).Return("", domain.ErrInvalidCredentials).Once()
			},
			wantStatusCode: http.StatusUnauthorized,
			wantBody:       `{"status":"error","code":401,"message":"invalid credentials"}` + "\n",
		},
		{
			name: "internal error",
			body: authHandler.LoginRequest{
				Login:    "test_login",
				Password: "test_password",
			},
			mockBehavior: func(repo *mocks.AuthService, body authHandler.LoginRequest) {
				repo.EXPECT().Login(mock.Anything, body.Login, body.Password).Return("", errors.New("internal error")).Once()
			},
			wantStatusCode: http.StatusInternalServerError,
			wantBody:       `{"status":"error","code":500,"message":"failed to login"}` + "\n",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			authService := mocks.NewAuthService(t)
			tc.mockBehavior(authService, tc.body)

			handler := authHandler.New(testutils.NewTestLogger(), authService)
			rec := httptest.NewRecorder()
			req := testutils.NewJSONRequest(t, http.MethodPost, "/auth/login", tc.body)

			handler.HandleLogin(rec, req)

			assert.Equal(t, tc.wantStatusCode, rec.Code)
			assert.Equal(t, tc.wantBody, rec.Body.String())
		})
	}
}
