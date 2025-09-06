package repository

import (
	"context"

	"github.com/ialekseychuk/my-place-auth/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

type repo struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *repo {
	return &repo{
		db: db,
	}
}


func (r *repo) GetByEmail(ctx context.Context, login string) (*domain.User, error) {
	var user domain.User
	err := r.db.QueryRow(ctx,
		`SELECT id, business_id, first_name, last_name, email, phone, password, role, is_active, created_at, updated_at
		 FROM users
		 WHERE email = $1`,
		login).Scan(&user.ID, &user.BusinessID, &user.FirstName, &user.LastName, &user.Email,
		&user.Phone, &user.Password, &user.Role, &user.IsActive, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

