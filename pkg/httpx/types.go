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
	Status  Status `json:"status"`
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type Middleware func(http.Handler) http.Handler

type AuthFunc func(token string) (context.Context, error)
