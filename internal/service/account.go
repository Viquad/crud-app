package service

import (
	"context"

	"github.com/Viquad/crud-app/internal/domain"
)

type AccountService struct {
	repo domain.AccountRepository
}

func NewAccountService(repo domain.AccountRepository) *AccountService {
	return &AccountService{
		repo: repo,
	}
}

func (b *AccountService) Create(ctx context.Context, account domain.Account) (*domain.Account, error) {
	return b.repo.Create(ctx, account)
}

func (b *AccountService) GetById(ctx context.Context, id int64) (*domain.Account, error) {
	return b.repo.GetById(ctx, id)
}

func (b *AccountService) List(ctx context.Context) ([]domain.Account, error) {
	return b.repo.List(ctx)
}

func (b *AccountService) UpdateById(ctx context.Context, id int64, inp domain.AccountUpdateInput) (*domain.Account, error) {
	return b.repo.UpdateById(ctx, id, inp)
}

func (b *AccountService) DeleteById(ctx context.Context, id int64) error {
	return b.repo.DeleteById(ctx, id)
}
