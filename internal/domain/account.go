package domain

import (
	"context"
	"time"
)

type Account struct {
	Id         int64     `form:"id" json:"id" example:"1"`
	UserId     int64     `form:"id" json:"user_id" example:"1"`
	Balance    int64     `form:"balance" json:"balance" example:"1000"`
	Currency   string    `form:"currency" json:"currency" example:"UAH"`
	LastUpdate time.Time `form:"lastUpdate" json:"lastUpdate" example:"2022-08-25T14:58:16.413065Z"`
}

type AccountCreateInput struct {
	Balance  int64  `form:"balance" json:"balance" binding:"required" example:"200"`
	Currency string `form:"currency" json:"currency" binding:"required" example:"UAH"`
}

type AccountUpdateInput struct {
	Balance *int64 `form:"balance" json:"balance" example:"1000"`
}

type AccountService interface {
	Create(ctx context.Context, inp AccountCreateInput) (*Account, error)
	List(ctx context.Context) ([]Account, error)
	GetById(ctx context.Context, id int64) (*Account, error)
	UpdateById(ctx context.Context, id int64, inp AccountUpdateInput) (*Account, error)
	DeleteById(ctx context.Context, id int64) error
}

type AccountRepository interface {
	Create(ctx context.Context, inp AccountCreateInput) (*Account, error)
	List(ctx context.Context) ([]Account, error)
	GetById(ctx context.Context, id int64) (*Account, error)
	UpdateById(ctx context.Context, id int64, inp AccountUpdateInput) (*Account, error)
	DeleteById(ctx context.Context, id int64) error
}
