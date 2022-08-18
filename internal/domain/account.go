package domain

import (
	"context"
	"time"
)

type Account struct {
	Id         int64     `form:"id" json:"id"`
	FirstName  string    `form:"firstName" json:"firstName" binding:"required"`
	LastName   string    `form:"lastName" json:"lastName" binding:"required"`
	Balance    int64     `form:"balance" json:"balance" binding:"required"`
	Currency   string    `form:"currency" json:"currency" binding:"required"`
	LastUpdate time.Time `form:"lastUpdate" json:"lastUpdate"`
}

type AccountUpdateInput struct {
	FirstName *string `form:"firstName" json:"firstName"`
	LastName  *string `form:"lastName" json:"lastName"`
	Balance   *int64  `form:"balance" json:"balance"`
	Currency  *string `form:"currency" json:"currency"`
}

type AccountService interface {
	Create(ctx context.Context, account Account) (*Account, error)
	All(ctx context.Context) ([]Account, error)
	GetById(ctx context.Context, id int64) (*Account, error)
	Update(ctx context.Context, id int64, inp AccountUpdateInput) (*Account, error)
	Delete(ctx context.Context, id int64) error
}

type AccountRepository interface {
	Create(ctx context.Context, account Account) (*Account, error)
	All(ctx context.Context) ([]Account, error)
	GetById(ctx context.Context, id int64) (*Account, error)
	Update(ctx context.Context, id int64, inp AccountUpdateInput) (*Account, error)
	Delete(ctx context.Context, id int64) error
}
