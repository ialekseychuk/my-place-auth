package handler

import (
	"context"
	"time"

	"github.com/ialekseychuk/my-place-identity/internal/domain"
	identityv1 "github.com/ialekseychuk/my-place-proto/gen/go/identity/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type IdentityHandler struct {
	identityv1.UnimplementedIdentityServer
	loginUC    domain.LoginUseCase
	registerUC domain.RegisterUseCase
	refreshUC  domain.RefreshUseCase
	validateUC domain.ValidateUseCase
	logoutUC   domain.LogoutUseCase
	getMeUC    domain.GetMeUseCase
}

func NewIdentityHandler(
	loginUC domain.LoginUseCase,
	registerUC domain.RegisterUseCase,
	refreshUC domain.RefreshUseCase,
	validateUC domain.ValidateUseCase,
	logoutUC domain.LogoutUseCase,
	getMeUC domain.GetMeUseCase,
) *IdentityHandler {
	return &IdentityHandler{
		loginUC:    loginUC,
		registerUC: registerUC,
		refreshUC:  refreshUC,
		validateUC: validateUC,
		logoutUC:   logoutUC,
		getMeUC:    getMeUC,
	}
}

// --------------------  RPC  --------------------

func (h *IdentityHandler) Login(ctx context.Context, req *identityv1.LoginRequest) (*identityv1.LoginResponse, error) {
	if req.Login == "" || req.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "login and password required")
	}
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	user, token, err := h.loginUC.Execute(ctx, req.Login, req.Password)
	if err != nil {
		return nil, handleError(err)
	}
	return &identityv1.LoginResponse{
		User:      mapUserToProto(user),
		AuthToken: mapTokenToProto(token),
	}, nil
}

func (h *IdentityHandler) Register(ctx context.Context, req *identityv1.RegisterRequest) (*identityv1.RegisterResponse, error) {
	if req.Email == "" || req.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "email and password required")
	}
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	domainReq := domain.RegisterRequest{
		Email:     req.Email,
		Phone:     req.Phone,
		Password:  req.Password,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Role:      req.Role,
	}
	user, token, err := h.registerUC.Execute(ctx, domainReq)
	if err != nil {
		return nil, handleError(err)
	}
	return &identityv1.RegisterResponse{
		User:      mapUserToProto(user),
		AuthToken: mapTokenToProto(token),
	}, nil
}

func (h *IdentityHandler) RefreshToken(ctx context.Context, req *identityv1.RefreshTokenRequest) (*identityv1.RefreshTokenResponse, error) {
	if req.RefreshToken == "" {
		return nil, status.Error(codes.InvalidArgument, "refresh token required")
	}
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	_, token, err := h.refreshUC.Execute(ctx, req.RefreshToken)
	if err != nil {
		return nil, handleError(err)
	}
	return &identityv1.RefreshTokenResponse{
		AuthToken: mapTokenToProto(token),
	}, nil
}

func (h *IdentityHandler) ValidateToken(ctx context.Context, req *identityv1.ValidateTokenRequest) (*identityv1.ValidateTokenResponse, error) {
	if req.AccessToken == "" {
		return nil, status.Error(codes.InvalidArgument, "access token required")
	}
	user, err := h.validateUC.Execute(ctx, req.AccessToken)
	if err != nil {
		return nil, handleError(err)
	}

	return &identityv1.ValidateTokenResponse{
		Valid: true,
		User:  mapUserToProto(user),
	}, nil
}

func (h *IdentityHandler) Logout(ctx context.Context, req *identityv1.LogoutRequest) (*identityv1.LogoutResponse, error) {
	if req.RefreshToken == "" {
		return nil, status.Error(codes.InvalidArgument, "refresh token required")
	}
	if err := h.logoutUC.Execute(ctx, req.RefreshToken); err != nil {
		return nil, handleError(err)
	}
	return &identityv1.LogoutResponse{}, nil
}

func (h *IdentityHandler) GetMe(ctx context.Context, _ *identityv1.GetMeRequest) (*identityv1.GetMeResponse, error) {
	userID, ok := ctx.Value("user_id").(string)
	if !ok || userID == "" {
		return nil, status.Error(codes.Unauthenticated, "user not found in context")
	}
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	user, err := h.getMeUC.Execute(ctx, userID)
	if err != nil {
		return nil, handleError(err)
	}
	return &identityv1.GetMeResponse{User: mapUserToProto(user)}, nil
}
