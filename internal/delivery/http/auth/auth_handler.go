package auth

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/SergeyBogomolovv/fitflow/internal/domain"
	"github.com/SergeyBogomolovv/fitflow/pkg/httpx"
	"github.com/go-playground/validator/v10"
)

type Service interface {
	Login(ctx context.Context, login, password string) (string, error)
}

type handler struct {
	logger   *slog.Logger
	validate *validator.Validate
	svc      Service
}

func New(logger *slog.Logger, svc Service) *handler {
	validate := validator.New(validator.WithRequiredStructEnabled())

	return &handler{logger, validate, svc}
}

func (h *handler) Handle(r *http.ServeMux) {
	r.HandleFunc("/auth/login", h.HandleLogin)
}

func (h *handler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	dto := new(LoginRequest)
	if err := httpx.DecodeBody(r, dto); err != nil {
		httpx.WriteError(w, "invalid body", http.StatusBadRequest)
		return
	}

	if err := h.validate.Struct(dto); err != nil {
		if err, ok := err.(validator.ValidationErrors); ok {
			httpx.WriteError(w, fmt.Sprintf("invalid request: %v", err), http.StatusBadRequest)
			return
		}
		httpx.WriteError(w, "failed to validate request", http.StatusInternalServerError)
		return
	}

	token, err := h.svc.Login(r.Context(), dto.Login, dto.Password)
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
