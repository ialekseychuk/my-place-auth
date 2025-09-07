package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/ialekseychuk/my-place-identity/internal/domain"
	"github.com/ialekseychuk/my-place-identity/internal/infrastructure"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type refreshUseCase struct {
	userRepo   domain.UserRepository
	tokenRepo  domain.TokenRepository
	jwt        *infrastructure.JWTManager
	accessTTL  time.Duration
	refreshTTL time.Duration
}

func NewRefresh(userRepo domain.UserRepository, tokenRepo domain.TokenRepository, jwt *infrastructure.JWTManager, access, refresh time.Duration) domain.RefreshUseCase {
	return &refreshUseCase{
		userRepo:   userRepo,
		tokenRepo:  tokenRepo,
		jwt:        jwt,
		accessTTL:  access,
		refreshTTL: refresh,
	}
}

func (r *refreshUseCase) Execute(ctx context.Context, refreshToken string) (*domain.User, *domain.AuthToken, error) {
	refreshHash := infrastructure.GenerateTokenHash(refreshToken)
	rt, err := r.tokenRepo.FindByHash(ctx, refreshHash)
	if err != nil {
		if errors.Is(err, domain.ErrRefreshTokenNotFound) {
			return nil, nil, status.Error(codes.Unauthenticated, "refresh token not found or expired")
		}
		return nil, nil, status.Error(codes.Internal, "failed to refresh token")
	}
	user, err := r.userRepo.GetByID(ctx, rt.UserID)
	if err != nil {
		return nil, nil, status.Error(codes.Internal, "failed to get user")
	}

	expireTime := time.Now().Add(r.accessTTL)

	accessToken, err := r.jwt.GenerateAccessToken(user, expireTime)
	if err != nil {
		return nil, nil, status.Error(codes.Internal, "failed to generate access token")
	}

	newRefreshRaw, err := infrastructure.GenerateRefreshToken()

	if err != nil {
		return nil, nil, status.Error(codes.Internal, "failed to generate refresh token")
	}
	newHash := infrastructure.GenerateTokenHash(newRefreshRaw)

	if err := r.tokenRepo.DeleteByHash(ctx, refreshHash); err != nil {
		return nil, nil, status.Error(codes.Internal, "failed to delete old refresh token")
	}

	if err := r.tokenRepo.Create(ctx, user.ID, newHash, time.Now().Add(r.refreshTTL)); err != nil {
		return nil, nil, status.Error(codes.Internal, "failed to create new refresh token")
	}
	return user, &domain.AuthToken{
		AccessToken:  accessToken,
		RefreshToken: newRefreshRaw,
		ExpiredAt:    expireTime,
	}, nil

}
