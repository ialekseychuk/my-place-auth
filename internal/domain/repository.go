package domain

import (
	"context"
	"time"
)

type UserRepository interface {
	Create(ctx context.Context, user *User) error
	GetByEmail(ctx context.Context, email string) (*User, error)
	GetByID(ctx context.Context, id string) (*User, error)
	//Update(ctx context.Context, user *domain.User) error
	//Delete(ctx context.Context, id string) error
}

type TokenRepository interface {
	Create(ctx context.Context, userID, tokenHash string, expiresAt time.Time) error
	FindByHash(ctx context.Context, hash string) (*RefreshToken, error)
	DeleteByHash(ctx context.Context, hash string) error
	DeleteAllByUserID(ctx context.Context, userID string) error
}
