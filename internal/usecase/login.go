package usecase

import (
	"context"
	"time"

	"github.com/ialekseychuk/my-place-identity/internal/domain"
	"github.com/ialekseychuk/my-place-identity/internal/infrastructure"
	"golang.org/x/crypto/bcrypt"
)

type loginUseCase struct {
	userRepo   domain.UserRepository
	tokenRepo  domain.TokenRepository
	jwt        *infrastructure.JWTManager
	accessTTL  time.Duration
	refreshTTL time.Duration
}

func NewLogin(userRepo domain.UserRepository, tokenRepo domain.TokenRepository,
	jwt *infrastructure.JWTManager, accessTTL, refreshTTL time.Duration) domain.LoginUseCase {
	return &loginUseCase{
		userRepo:   userRepo,
		tokenRepo:  tokenRepo,
		jwt:        jwt,
		accessTTL:  accessTTL,
		refreshTTL: refreshTTL,
	}
}

func (u *loginUseCase) Execute(ctx context.Context, email, password string) (*domain.User, *domain.AuthToken, error) {
	user, err := u.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, nil, domain.ErrInvalidCredentials
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, nil, domain.ErrInvalidCredentials
	}

	accessExp := time.Now().Add(u.accessTTL)
	accessToken, err := u.jwt.GenerateAccessToken(user, accessExp)

	if err != nil {
		return nil, nil, err
	}

	refreshRaw, _ := infrastructure.GenerateRefreshToken()
	refreshHash := infrastructure.GenerateTokenHash(refreshRaw)

	if err := u.tokenRepo.Create(ctx, user.ID, refreshHash, time.Now().Add(u.refreshTTL)); err != nil {
		return nil, nil, err
	}

	return user, &domain.AuthToken{
		AccessToken:  accessToken,
		RefreshToken: refreshRaw,
		ExpiredAt:    accessExp,
		TokenType:    "Bearer",
	}, nil
}
