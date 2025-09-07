package handler

import (
	"errors"

	"github.com/ialekseychuk/my-place-identity/internal/domain"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func handleError(err error) error {
	switch {
	case errors.Is(err, domain.ErrInvalidCredentials):
		return status.Error(codes.Unauthenticated, "invalid credentials")
	case errors.Is(err, domain.ErrEmailExists):
		return status.Error(codes.AlreadyExists, "email already registered")
	case errors.Is(err, domain.ErrRefreshTokenNotFound):
		return status.Error(codes.Unauthenticated, "refresh token not found or expired")
	default:
		return status.Error(codes.Internal, "internal error")
	}
}