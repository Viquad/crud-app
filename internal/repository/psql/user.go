package psql

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Viquad/crud-app/internal/domain"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) Create(ctx context.Context, input domain.SignUpInput) error {
	query := "INSERT INTO users (first_name, last_name, email, password) VALUES ($1, $2, $3, $4)"
	_, err := r.db.ExecContext(ctx, query, input.FirstName, input.LastName, input.Email, input.Password)

	return err
}

func (r *UserRepository) GetByCredentials(ctx context.Context, input domain.SignInInput) (*domain.User, error) {
	var user domain.User
	query := "SELECT id, first_name, last_name, email, password, registered_at FROM users WHERE email=$1 AND password=$2"
	err := r.db.QueryRowContext(ctx, query, input.Email, input.Password).
		Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.Password, &user.RegisteredAt)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrUserNotFound
	}

	return &user, err
}
