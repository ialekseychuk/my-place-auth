package domain

import "errors"

var (
	ErrInvalidCredentials   = errors.New("invalid credentials")
	ErrEmailExists          = errors.New("email already registered")
	ErrRefreshTokenNotFound = errors.New("refresh token not found or expired")
	ErrTokenExpired         = errors.New("token expired")
	ErrTokenMalformed       = errors.New("token malformed")
	ErrUserNotActive        = errors.New("user not active")
)
