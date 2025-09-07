package repository

import (
	"context"

	"github.com/ialekseychuk/my-place-identity/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

type userRepo struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *userRepo {
	return &userRepo{
		db: db,
	}
}

func (r *userRepo) GetByEmail(ctx context.Context, login string) (*domain.User, error) {
	var user domain.User
	err := r.db.QueryRow(ctx,
		`SELECT id, first_name, last_name, email, phone, password, role, is_active, email_verified, phone_verified created_at, updated_at
		 FROM users
		 WHERE email = $1`,
		login).Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email,
		&user.Phone, &user.Password, &user.Role, &user.IsActive, &user.EmailVerified, 
		&user.PhoneVerified, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepo) Create(ctx context.Context, user *domain.User) error {
	_, err := r.db.Exec(ctx,
		`INSERT INTO users (id, first_name, last_name, email, phone, password, role, 
		is_active, email_verified, phone_verified, user created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`,
		user.ID, user.FirstName, user.LastName, user.Email, user.Phone,
		user.Password, user.Role, user.IsActive, 
		user.EmailVerified, user.PhoneVerified, user.CreatedAt, user.UpdatedAt)
	return err
}