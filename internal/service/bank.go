package service

import (
	"context"

	"github.com/Viquad/crud-app/internal/domain"
)

type BankRepository interface {
	CreateAccount(ctx context.Context, account domain.Account) error
	ReadAccount(ctx context.Context, id int64) (domain.Account, error)
	ReadAllAccounts(ctx context.Context) ([]domain.Account, error)
	UpdateAccount(ctx context.Context, id int64, inp domain.AccountUpdateInput) error
	DeleteAccount(ctx context.Context, id int64) error
}

type Bank struct {
	repo BankRepository
}

func NewBank(repo BankRepository) *Bank {
	return &Bank{
		repo: repo,
	}
}

func (b *Bank) CreateAccount(ctx context.Context, account domain.Account) error {
	return b.repo.CreateAccount(ctx, account)
}

func (b *Bank) ReadAccount(ctx context.Context, id int64) (domain.Account, error) {
	return b.repo.ReadAccount(ctx, id)
}

func (b *Bank) ReadAllAccounts(ctx context.Context) ([]domain.Account, error) {
	return b.repo.ReadAllAccounts(ctx)
}

func (b *Bank) UpdateAccount(ctx context.Context, id int64, inp domain.AccountUpdateInput) error {
	return b.repo.UpdateAccount(ctx, id, inp)
}

func (b *Bank) DeleteAccount(ctx context.Context, id int64) error {
	return b.repo.DeleteAccount(ctx, id)
}
