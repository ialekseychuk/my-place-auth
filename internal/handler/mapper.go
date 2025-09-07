package handler

import (
	identityv1 "github.com/ialekseychuk/my-place-proto/gen/go/identity/v1"
	"github.com/ialekseychuk/my-place-identity/internal/domain"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func mapUserToProto(u *domain.User) *identityv1.User {
	return &identityv1.User{
		Id:            u.ID,
		Email:         u.Email,
		Phone:         u.Phone,
		FirstName:     u.FirstName,
		LastName:      u.LastName,
		Role:          u.Role,
		IsActive:      u.IsActive,
		EmailVerified: u.EmailVerified,
		PhoneVerified: u.PhoneVerified,
		CreatedAt:     timestamppb.New(u.CreatedAt),
		UpdatedAt:     timestamppb.New(u.UpdatedAt),
	}
}

func mapTokenToProto(t *domain.AuthToken) *identityv1.AuthToken {
	return &identityv1.AuthToken{
		AccessToken:  t.AccessToken,
		RefreshToken: t.RefreshToken,
		ExpiredAt:    timestamppb.New(t.ExpiredAt),
		TokenType:    t.TokenType,
	}
}