package domain

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
)

type keyType string

const UserIdKey keyType = "user_id"

type User struct {
	Id           int64     `form:"id" json:"id" example:"1"`
	FirstName    string    `form:"firstName" json:"firstName" binding:"required"`
	LastName     string    `form:"lastName" json:"lastName" binding:"required"`
	Email        string    `form:"email" json:"email" binding:"required"`
	Password     string    `form:"password" json:"password" binding:"required"`
	RegisteredAt time.Time `form:"lastUpdate" json:"lastUpdate"`
}

type SignUpInput struct {
	FirstName string `form:"firstName" json:"firstName" binding:"required,gte=2" example:"Oleksii"`
	LastName  string `form:"lastName" json:"lastName" binding:"required,gte=2" example:"Filatov"`
	Email     string `form:"email" json:"email" binding:"required,email" example:"ofilatov@gmail.com"`
	Password  string `form:"password" json:"password" binding:"required,gte=8" example:"TheBestGuy99"`
}

type SignInInput struct {
	Email    string `form:"email" json:"email" binding:"required,email" example:"ofilatov@gmail.com"`
	Password string `form:"password" json:"password" binding:"required,gte=8" example:"TheBestGuy99"`
}

type UserService interface {
	Create(ctx context.Context, input SignUpInput) error
	GetTokenByCredentials(ctx context.Context, input SignInInput) (string, string, error)
	ParseToken(ctx context.Context, token string) (int64, error)
	RefreshTokens(ctx context.Context, token string) (string, string, error)
	InitSession(c *gin.Context, input SignInInput) error
	GetSession(c *gin.Context) (int64, error)
	DropSession(c *gin.Context) error
}

type UserRepository interface {
	Create(ctx context.Context, input SignUpInput) error
	GetByCredentials(ctx context.Context, input SignInInput) (*User, error)
}
