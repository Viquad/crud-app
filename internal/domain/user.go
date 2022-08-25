package domain

import (
	"context"
	"time"
)

type User struct {
	Id           int64     `form:"id" json:"id"`
	FirstName    string    `form:"firstName" json:"firstName" binding:"required"`
	LastName     string    `form:"lastName" json:"lastName" binding:"required"`
	Email        string    `form:"email" json:"email" binding:"required"`
	Password     string    `form:"password" json:"password" binding:"required"`
	RegisteredAt time.Time `form:"lastUpdate" json:"lastUpdate"`
}

type SignUpInput struct {
	FirstName string `form:"firstName" json:"firstName" binding:"required,gte=2"`
	LastName  string `form:"lastName" json:"lastName" binding:"required,gte=2"`
	Email     string `form:"email" json:"email" binding:"required,email"`
	Password  string `form:"password" json:"password" binding:"required,gte=8"`
}

type SignInInput struct {
	Email    string `form:"email" json:"email" binding:"required,email"`
	Password string `form:"password" json:"password" binding:"required,gte=8"`
}

type UserService interface {
	Create(ctx context.Context, input SignUpInput) error
	GetTokenByCredentials(ctx context.Context, input SignInInput) (string, error)
	ParseToken(ctx context.Context, token string) (int64, error)
}

type UserRepository interface {
	Create(ctx context.Context, input SignUpInput) error
	GetByCredentials(ctx context.Context, input SignInInput) (*User, error)
}
