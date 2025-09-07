package usecase

import (
	"context"
	"time"

	"github.com/ialekseychuk/my-place-identity/internal/domain"
	"github.com/ialekseychuk/my-place-identity/internal/infrastructure"
	"github.com/jackc/pgx/v5"
)

type registerUC struct {
	userRepo   domain.UserRepository
	tokenRepo  domain.TokenRepository
	jwt        *infrastructure.JWTManager
	accessTTL  time.Duration
	refreshTTL time.Duration
}

func NewRegister(userRepo domain.UserRepository, tokenRepo domain.TokenRepository,
	jwt *infrastructure.JWTManager, accessTTL, refreshTTL time.Duration) domain.RegisterUseCase {
	return &registerUC{
		userRepo:   userRepo,
		tokenRepo:  tokenRepo,
		jwt:        jwt,
		accessTTL:  accessTTL,
		refreshTTL: refreshTTL,
	}
}

func (r *registerUC) Execute(ctx context.Context, req domain.RegisterRequest) (*domain.User, *domain.AuthToken, error) {
	_, err := r.userRepo.GetByEmail(ctx, req.Email)

	if err != nil && err != pgx.ErrNoRows {
		return nil, nil, err
	}

	hash, err := infrastructure.GeneratePassworHash(req.Password)
	if err != nil {
		return nil, nil, err
	}
	user := &domain.User{
		Email:         req.Email,
		Password:      string(hash),
		FirstName:     req.FirstName,
		LastName:      req.LastName,
		Role:          req.Role,
		Phone:         req.Phone,
		IsActive:      true,
		EmailVerified: false,
		PhoneVerified: false,
	}

	if err := r.userRepo.Create(ctx, user); err != nil {
		return nil, nil, err
	}

	accessExp := time.Now().Add(r.accessTTL)
	accessToken, err := r.jwt.GenerateAccessToken(user, accessExp)

	if err != nil {
		return nil, nil, err
	}

	refreshRaw, _ := infrastructure.GenerateRefreshToken()
	refreshHash := infrastructure.GenerateTokenHash(refreshRaw)

	if err := r.tokenRepo.Create(ctx, user.ID, refreshHash, time.Now().Add(r.refreshTTL)); err != nil {
		return nil, nil, err
	}

	return user, &domain.AuthToken{
		AccessToken:  accessToken,
		RefreshToken: refreshRaw,
		ExpiredAt:    accessExp,
		TokenType:    "Bearer",
	}, nil

}
