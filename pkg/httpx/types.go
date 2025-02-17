package httpx

import (
	"context"
	"net/http"
)

type Status string

const (
	StatusSuccess Status = "success"
	StatusError   Status = "error"
)

type Response struct {
	Status  Status `json:"status" example:"success"`
	Code    int    `json:"code" example:"200"`
	Message string `json:"message" example:"Запрос выполнен успешно"`
}

type Middleware func(http.Handler) http.Handler

type AuthFunc func(token string) (context.Context, error)
