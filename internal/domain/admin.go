package domain

import "errors"

type Admin struct {
	Login    string
	Password []byte
}

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrAdminAlreadyExists = errors.New("admin already exists")
	ErrAdminNotFound      = errors.New("admin not found")
)
