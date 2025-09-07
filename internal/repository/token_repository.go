package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/ialekseychuk/my-place-identity/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type tokenRepo struct {
	db *pgxpool.Pool
}

func NewTokenRepository(db *pgxpool.Pool) *tokenRepo {
	return &tokenRepo{
		db: db,
	}
}

func (r *tokenRepo) Create(ctx context.Context, userID, tokenHash string, expiresAt time.Time) error {
	const query = `
		INSERT INTO refresh_tokens (user_id, token_hash, expires_at)
		VALUES ($1, $2, $3)
		ON CONFLICT (token_hash) DO UPDATE
		   SET expires_at = EXCLUDED.expires_at,
		       created_at = now()
	`
	_, err := r.db.Exec(ctx, query, userID, tokenHash, expiresAt)
	if err != nil {
		return fmt.Errorf("save refresh token: %w", err)
	}
	return nil
}
func (r *tokenRepo) FindByHash(ctx context.Context, hash string) (*domain.RefreshToken, error) {
	const query = `
	SELECT id, user_id, token_hash, expires_at, created_at
	FROM refresh_tokens
	Where token_hash = $1
	AND expires_at > now()
	LIMIT 1`
	var rt domain.RefreshToken
	err := r.db.QueryRow(ctx, query, hash).Scan(
		&rt.ID,
		&rt.UserID,
		&rt.TokenHash,
		&rt.ExpiresAt,
		&rt.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrRefreshTokenNotFound
		}
		return nil, fmt.Errorf("failed to find refresh token: %w", err)
	}
	return &rt, nil
}
func (r *tokenRepo) DeleteByHash(ctx context.Context, hash string) error {
	const query = `
	DELETE FROM refresh_tokens
	Where token_hash = $1`
	tag, err := r.db.Exec(ctx, query, hash)
	if err != nil {
		return fmt.Errorf("failed to delete refresh token: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrRefreshTokenNotFound
	}
	return nil
}
func (r *tokenRepo) DeleteAllByUserID(ctx context.Context, userID string) error {
	const query = `
	DELETE FROM refresh_tokens
	Where user_id = $1`
	_, err := r.db.Exec(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("failed to delete refresh token: %w", err)

	}
	return nil
}
