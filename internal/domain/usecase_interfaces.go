package domain

import "context"

type LoginUseCase interface {
	Execute(ctx context.Context, email, password string) (*User, *AuthToken, error)
}

type RegisterUseCase interface {
	Execute(ctx context.Context, req RegisterRequest) (*User, *AuthToken, error)
}

type RefreshUseCase interface {
	Execute(ctx context.Context, refreshToken string) (*User, *AuthToken, error)
}

type ValidateUseCase interface {
	Execute(ctx context.Context, accessToken string) (*User, error)
}

type LogoutUseCase interface {
	Execute(ctx context.Context, refreshToken string) error
}

type GetMeUseCase interface {
	Execute(ctx context.Context, userID string) (*User, error)
}
