package infrastructure

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/ialekseychuk/my-place-identity/internal/domain"
)

type JWTManager struct {
	secret string
}

func NewJWTManager(secret string) *JWTManager {
	return &JWTManager{secret: secret}
}

func (j *JWTManager) GenerateAccessToken(user *domain.User, expiry time.Time) (string, error) {
	claims := jwt.MapClaims{
		"sub":  user.ID,
		"email": user.Email,
		"role": user.Role,
		"exp":  expiry.Unix(),
		"iat":  time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.secret))
}

func (j *JWTManager) ValidateAccessToken(tokenStr string) (*domain.TokenClaims, error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, domain.ErrTokenMalformed
		}
		return []byte(j.secret), nil
	})
	if err != nil || !token.Valid {
		return nil, domain.ErrTokenExpired
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, domain.ErrTokenMalformed
	}
	return &domain.TokenClaims{
		UserID: claims["sub"].(string),
		Email:  claims["email"].(string),
		Role:   claims["role"].(string),
	}, nil
}