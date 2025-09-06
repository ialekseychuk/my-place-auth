package domain

import (
	"context"

	
)

type UserRepository interface {
	//Create(ctx context.Context, user *domain.User) error
	GetByEmail(ctx context.Context, email string) (*User, error)
	//GetByID(ctx context.Context, id string) (*domain.User, error)
	//Update(ctx context.Context, user *domain.User) error
	//Delete(ctx context.Context, id string) error
}