package auth

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/ialekseychuk/my-place-auth/internal/domain"
	authv1 "github.com/ialekseychuk/my-place-proto/gen/go/auth/v1"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type service struct {
	ur domain.UserRepository
	jwtSecret              string
	accessTokenExpiryTime  time.Duration
	refreshTokenExpiryTime time.Duration
	authv1.UnimplementedAuthServer
}

func NewAuthService(repo domain.UserRepository, jwtSecret string) *service {
	return &service{
		ur: repo,
		jwtSecret:              jwtSecret,
		accessTokenExpiryTime:  24 * time.Hour,      // 24 hours
		refreshTokenExpiryTime: 30 * 24 * time.Hour, // 30 days
	}
}

func (s *service) Login(ctx context.Context, req *authv1.LoginRequest) (*authv1.LoginResponse, error) {
	user, err := s.ur.GetByEmail(ctx, req.Login)
	
	if err != nil {
		return nil, status.Error(codes.NotFound, "user not found")
	}
	if user.IsActive == false {
		return nil, status.Error(codes.PermissionDenied, "user is not active")
	}

	if !s.verifyPassword(user.Password, req.Password) {
		return nil, status.Error(codes.Unauthenticated, "invalid credentials")
	}
	accessToken, expired_at, err := s.generateAccessToken(user)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to generate access token: %v", err)
	}

	refreshToken, err := s.generateRefreshToken()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to generate refresh token: %v", err)
	}
	resp := &authv1.LoginResponse{
		User:  &authv1.User{
			Id: user.ID, 
			BusinessId: user.BusinessID, 
			FirstName: user.FirstName, 
			LastName: user.LastName, 
			Email: user.Email, 
			Phone: user.Phone, 
			Role: user.Role, 
			IsActive: user.IsActive, 
			CreatedAt: user.CreatedAt.Format(time.RFC3339), 
			UpdatedAt: user.UpdatedAt.Format(time.RFC3339),
		},
		AuthToken : &authv1.AuthToken{
			AccessToken: accessToken,
			RefreshToken: refreshToken,
			ExpiredAt: timestamppb.New(expired_at),
			TokenType: "Bearer",
		},
	}
	return resp, nil
}

func (s *service) verifyPassword(storedPassword, providedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(providedPassword))
	return err == nil
}

func (s *service) generateAccessToken(user *domain.User) (string, time.Time, error) {
	now := time.Now()
	expiryTime := now.Add(s.accessTokenExpiryTime)

	claims := &domain.JWTClaims{
		UserID:     user.ID,
		BusinessID: user.BusinessID,
		Email:      user.Email,
		Role:       user.Role,
		IssuedAt:   now.Unix(),
		ExpiresAt:  expiryTime.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", time.Time{}, fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, expiryTime, nil
}

func (s *service) generateRefreshToken() (string, error) {
	bytes := make([]byte, 32)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", fmt.Errorf("failed to generate refresh token: %w", err)
	}
	return hex.EncodeToString(bytes), nil
}
