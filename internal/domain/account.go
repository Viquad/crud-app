package domain

import (
	"context"
	"time"
)

type Account struct {
	Id         int64     `json:"id"`
	FirstName  string    `json:"firstName"`
	LastName   string    `json:"lastName"`
	Balance    int64     `json:"balance"`
	Currency   string    `json:"currency"`
	LastUpdate time.Time `json:"lastUpdate"`
}

type AccountUpdateInput struct {
	FirstName *string `json:"firstName"`
	LastName  *string `json:"lastName"`
	Balance   *int64  `json:"balance"`
	Currency  *string `json:"currency"`
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
