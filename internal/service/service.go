package service

import (
	"github.com/Viquad/crud-app/internal/domain"
)

type Repositories interface {
	GetAccountRepository() domain.AccountRepository
}

type Services struct {
	accountService *AccountService
}

func (ss *Services) GetAccountService() domain.AccountService {
	return ss.accountService
}

func NewServices(repo Repositories) *Services {
	return &Services{
		accountService: NewAccountService(repo.GetAccountRepository()),
	}
}
