package usecase

import (
	"context"

	"github.com/ialekseychuk/my-place-identity/internal/domain"
	"github.com/ialekseychuk/my-place-identity/internal/infrastructure"
)

type logoutUseCase struct {
	tokenRepo domain.TokenRepository
}

func NewLogout(tokenRepo domain.TokenRepository) *logoutUseCase {
	return &logoutUseCase{
		tokenRepo: tokenRepo,
	}
}

func (u *logoutUseCase) Execute(ctx context.Context, refreshToken string) error {
	hash := infrastructure.GenerateTokenHash(refreshToken)
	return u.tokenRepo.DeleteByHash(ctx, hash)
}
