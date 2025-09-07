package repository

import (
	"context"
	"errors"

	"github.com/ialekseychuk/my-place-identity/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
)

type userRepo struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *userRepo {
	return &userRepo{
		db: db,
	}
}

func (r *userRepo) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	const sql = `SELECT id, first_name, last_name, email, phone, password, role, is_active, email_verified, phone_verified, created_at, updated_at
		 FROM users
		 WHERE email = $1`

	var user domain.User
	err := r.db.QueryRow(ctx,
		sql,
		email).Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email,
		&user.Phone, &user.Password, &user.Role, &user.IsActive, &user.EmailVerified,
		&user.PhoneVerified, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		logrus.Debugf("Get user by email: %s, error: %s", email, err)
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepo) Create(ctx context.Context, user *domain.User) error {
	const sql = `INSERT INTO users (first_name, last_name, email, phone, password, role, 
			is_active, email_verified, phone_verified, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		 RETURNING id`
	err := r.db.QueryRow(ctx, sql,
		user.FirstName, user.LastName, user.Email, user.Phone,
		user.Password, user.Role, user.IsActive,
		user.EmailVerified, user.PhoneVerified, user.CreatedAt, user.UpdatedAt).Scan(&user.ID)
	return err
}

func (r *userRepo) GetByID(ctx context.Context, id string) (*domain.User, error) {
	var user domain.User
	err := r.db.QueryRow(ctx,
		`SELECT id, first_name, last_name, email, phone, password, role, is_active, email_verified, phone_verified, created_at, updated_at
		 FROM users
		 WHERE id = $1`,
		id).Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email,
		&user.Phone, &user.Password, &user.Role, &user.IsActive, &user.EmailVerified,
		&user.PhoneVerified, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}
