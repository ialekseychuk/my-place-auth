package grpc

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/ialekseychuk/my-place-identity/internal/domain"
	identityv1 "github.com/ialekseychuk/my-place-proto/gen/go/identity/v1"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type IdentityHandler struct {
	identityv1.UnimplementedAuthServer
	userRepo               domain.UserRepository
	jwtSecret              string
	accessTokenExpiryTime  time.Duration
	refreshTokenExpiryTime time.Duration
}

func NewIdentityHandler(userRepo domain.UserRepository, jwtSecret string) *IdentityHandler {
	return &IdentityHandler{
		userRepo:               userRepo,
		jwtSecret:              jwtSecret,
		accessTokenExpiryTime:  24 * time.Hour,      // 24 hours
		refreshTokenExpiryTime: 30 * 24 * time.Hour, // 30 days
	}
}

func (h *IdentityHandler) Login(ctx context.Context, req *identityv1.LoginRequest) (*identityv1.LoginResponse, error) {

	if req.Login == "" || req.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "login and password are required")
	}

	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	user, err := h.userRepo.GetByEmail(ctx, req.Login)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to get user by email")
	}

	if !h.verifyPassword(user.Password, req.Password) {
		return nil, status.Error(codes.Unauthenticated, "invalid credentials")
	}

	token, expired_at, err := h.generateAccessToken(user)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to generate access token")
	}
	refreshToken, err := h.generateRefreshToken()
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to generate refresh token")
	}

	return &identityv1.LoginResponse{
		User: &identityv1.User{
			Id:        user.ID,
			Email:     user.Email,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Phone:     user.Phone,
			Role:      user.Role,
			IsActive:  user.IsActive,
			EmailVerified: user.EmailVerified,
			CreatedAt: timestamppb.New(user.CreatedAt),
			UpdatedAt: timestamppb.New(user.UpdatedAt),
		},
		AuthToken: &identityv1.AuthToken{
			AccessToken:  token,
			RefreshToken: refreshToken,
			ExpiredAt:    timestamppb.New(expired_at),
		},
	}, nil
}


func (h *IdentityHandler) Register(ctx context.Context, req *identityv1.RegisterRequest) (*identityv1.RegisterResponse, error) {
	usr := &domain.User{
		FirstName:     req.FirstName,
		LastName:      req.LastName,
		Email:         req.Email,
		Phone:         req.Phone,
		Password:      req.Password,
		Role:          req.Role,
		IsActive:      true,
		EmailVerified: false,
		PhoneVerified: false,
	}
	h.userRepo.
}
func (h *IdentityHandler) Logout(context.Context, *identityv1.LogoutRequest) (*identityv1.LogoutResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Logout not implemented")
}
func (h *IdentityHandler) RefreshToken(context.Context, *identityv1.RefreshTokenRequest) (*identityv1.RefreshTokenResponse, error) {

	return nil, status.Errorf(codes.Unimplemented, "method RefreshToken not implemented")
}
func (h *IdentityHandler) ValidateToken(context.Context, *identityv1.ValidateTokenRequest) (*identityv1.ValidateTokenResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ValidateToken not implemented")
}
func (h *IdentityHandler) GetMe(context.Context, *identityv1.GetMeRequest) (*identityv1.GetMeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetMe not implemented")
}



func (h *IdentityHandler) verifyPassword(storedPassword, providedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(providedPassword))
	return err == nil
}

func (h *IdentityHandler) generateAccessToken(user *domain.User) (string, time.Time, error) {
	now := time.Now()
	expiryTime := now.Add(h.accessTokenExpiryTime)

	claims := &domain.JWTClaims{
		UserID:     user.ID,
		BusinessID: user.BusinessID,
		Email:      user.Email,
		Role:       user.Role,
		IssuedAt:   now.Unix(),
		ExpiresAt:  expiryTime.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(h.jwtSecret))
	if err != nil {
		return "", time.Time{}, fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, expiryTime, nil
}

func (s *IdentityHandler) generateRefreshToken() (string, error) {
	bytes := make([]byte, 32)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", fmt.Errorf("failed to generate refresh token: %w", err)
	}
	return hex.EncodeToString(bytes), nil
}
