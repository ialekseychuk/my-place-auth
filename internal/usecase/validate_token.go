package usecase

import (
	"context"

	"github.com/ialekseychuk/my-place-identity/internal/domain"
	"github.com/ialekseychuk/my-place-identity/internal/infrastructure"
)

type validateTokenUseCase struct {
	userRepo domain.UserRepository
	jwt      *infrastructure.JWTManager
}

func NewValidateToken(userRepo domain.UserRepository, jwt *infrastructure.JWTManager) domain.ValidateUseCase {
	return &validateTokenUseCase{
		userRepo: userRepo,
		jwt:      jwt,
	}
}

func (u *validateTokenUseCase) Execute(ctx context.Context, accessToken string) (*domain.User, error) {
	claims, err := u.jwt.ValidateAccessToken(accessToken)
	if err != nil {
		return nil, domain.ErrTokenExpired
	}

	user, err := u.userRepo.GetByID(ctx, claims.UserID)
	if err != nil {
		return nil, err
	}
	if !user.IsActive {
		return nil, domain.ErrUserNotActive
	}
	return user, nil
}
