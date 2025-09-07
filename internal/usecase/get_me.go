package usecase

import (
	"context"

	"github.com/ialekseychuk/my-place-identity/internal/domain"
)

type getMeUseCase struct {
	userRepo  domain.UserRepository
}

func NewGetMe(userRepo domain.UserRepository) domain.GetMeUseCase {
	return &getMeUseCase{
		userRepo: userRepo,
	}
}

func (u *getMeUseCase) Execute(ctx context.Context, userID string) (*domain.User, error){
	return u.userRepo.GetByID(ctx, userID)
}