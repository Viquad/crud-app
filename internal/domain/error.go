package domain

import "errors"

var (
	ErrNotExist          = errors.New("row does not exist")
	ErrUpdateFailed      = errors.New("update failed")
	ErrDeleteFailed      = errors.New("delete failed")
	ErrInvalidId         = errors.New("invalid id")
	ErrUserNotFound      = errors.New("user with such credentials not found")
	ErrInvalidClaims     = errors.New("invalid claims")
	ErrInvalidToken      = errors.New("invalid token")
	ErrUserAlreadyExists = errors.New("user already exists")
)
