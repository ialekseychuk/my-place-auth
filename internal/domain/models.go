package domain

import "time"

type User struct {
	ID            string
	Email         string
	Phone         string
	Password      string // bcrypt hash
	FirstName     string
	LastName      string
	Role          string
	IsActive      bool
	EmailVerified bool
	PhoneVerified bool
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type AuthToken struct {
	AccessToken  string
	RefreshToken string
	ExpiredAt    time.Time
	TokenType    string
}

type RegisterRequest struct {
	Email     string
	Phone     string
	Password  string // plaintext
	FirstName string
	LastName  string
	Role      string
}

type RefreshTokenClaims struct {
	UserID string
}

type RefreshToken struct {
	ID        string
	UserID    string
	TokenHash string
	ExpiresAt time.Time
	CreatedAt time.Time
}

type TokenClaims struct {
	UserID string
	Email  string
	Role   string
}
