package auth

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/SergeyBogomolovv/fitflow/internal/domain"
	"github.com/SergeyBogomolovv/fitflow/pkg/httpx"
	"github.com/go-playground/validator/v10"
)

type AuthService interface {
	Login(ctx context.Context, login, password string) (string, error)
}

type handler struct {
	logger   *slog.Logger
	validate *validator.Validate
	authSvc  AuthService
}

func New(logger *slog.Logger, authSvc AuthService) *handler {
	validate := validator.New(validator.WithRequiredStructEnabled())

	return &handler{logger, validate, authSvc}
}

func (h *handler) Init(r *http.ServeMux) {
	r.HandleFunc("POST /auth/login", h.HandleLogin)
}

// @Summary      Вход в учетную запись администратора
// @Description  Учетные записи администратора создаются через cli утилиту
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        input  body      LoginRequest  true  "Данные для входа"
// @Success      200    {object}  LoginResponse
// @Failure      400    {object}  httpx.Response  "Неверный формат данных"
// @Failure      401    {object}  httpx.Response  "Неверные данные для входа"
// @Failure      500    {object}  httpx.Response  "Внутренняя ошибка сервера"
// @Router       /auth/login [post]
func (h *handler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	var dto LoginRequest
	if err := httpx.DecodeBody(r, &dto); err != nil {
		httpx.WriteError(w, "invalid body", http.StatusBadRequest)
		return
	}

	if err := h.validate.Struct(dto); err != nil {
		httpx.WriteError(w, "invalid payload", http.StatusBadRequest)
		return
	}

	token, err := h.authSvc.Login(r.Context(), dto.Login, dto.Password)
	if err != nil {
		if errors.Is(err, domain.ErrInvalidCredentials) {
			httpx.WriteError(w, "invalid credentials", http.StatusUnauthorized)
			return
		}
		httpx.WriteError(w, "failed to login", http.StatusInternalServerError)
		return
	}

	httpx.WriteJSON(w, LoginResponse{Token: token}, http.StatusOK)
}
