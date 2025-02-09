package auth_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	authHandler "github.com/SergeyBogomolovv/fitflow/internal/delivery/http/auth"
	"github.com/SergeyBogomolovv/fitflow/internal/domain"
	testutils "github.com/SergeyBogomolovv/fitflow/pkg/test_utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestAuthHandler_Login(t *testing.T) {
	mockSvc := new(authServiceMock)
	handler := authHandler.New(testutils.NewTestLogger(), mockSvc)

	t.Run("success", func(t *testing.T) {
		reqBody := authHandler.LoginRequest{
			Login:    "test_user",
			Password: "correct_password",
		}
		mockSvc.On("Login", mock.Anything, reqBody.Login, reqBody.Password).
			Return("valid_token", nil).Once()

		rec := httptest.NewRecorder()
		req := testutils.NewJSONRequest(t, http.MethodPost, "/auth/login", reqBody)

		handler.HandleLogin(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var resp authHandler.LoginResponse
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
		assert.Equal(t, "valid_token", resp.Token)
		mockSvc.AssertExpectations(t)
	})

	t.Run("invalid payload", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBufferString("invalid_json"))

		handler.HandleLogin(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Contains(t, rec.Body.String(), "invalid body")
	})

	t.Run("wrong password", func(t *testing.T) {
		reqBody := authHandler.LoginRequest{
			Login:    "test_user",
			Password: "wrong_password",
		}
		mockSvc.On("Login", mock.Anything, reqBody.Login, reqBody.Password).
			Return("", domain.ErrInvalidCredentials).Once()

		rec := httptest.NewRecorder()
		req := testutils.NewJSONRequest(t, http.MethodPost, "/auth/login", reqBody)

		handler.HandleLogin(rec, req)

		assert.Equal(t, http.StatusUnauthorized, rec.Code)
		assert.Contains(t, rec.Body.String(), "invalid credentials")
		mockSvc.AssertExpectations(t)
	})

	t.Run("wrong login", func(t *testing.T) {
		reqBody := authHandler.LoginRequest{
			Login:    "wrong_user",
			Password: "some_password",
		}
		mockSvc.On("Login", mock.Anything, reqBody.Login, reqBody.Password).
			Return("", domain.ErrInvalidCredentials).Once()

		rec := httptest.NewRecorder()
		req := testutils.NewJSONRequest(t, http.MethodPost, "/auth/login", reqBody)

		handler.HandleLogin(rec, req)

		assert.Equal(t, http.StatusUnauthorized, rec.Code)
		assert.Contains(t, rec.Body.String(), "invalid credentials")
		mockSvc.AssertExpectations(t)
	})

	t.Run("error", func(t *testing.T) {
		reqBody := authHandler.LoginRequest{
			Login:    "test_user",
			Password: "some_password",
		}
		mockSvc.On("Login", mock.Anything, reqBody.Login, reqBody.Password).
			Return("", errors.New("unexpected error")).Once()

		rec := httptest.NewRecorder()
		req := testutils.NewJSONRequest(t, http.MethodPost, "/auth/login", reqBody)

		handler.HandleLogin(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.Contains(t, rec.Body.String(), "failed to login")
		mockSvc.AssertExpectations(t)
	})
}

type authServiceMock struct {
	mock.Mock
}

func (m *authServiceMock) Login(ctx context.Context, login, password string) (string, error) {
	args := m.Called(ctx, login, password)
	return args.String(0), args.Error(1)
}
