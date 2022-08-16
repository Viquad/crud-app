package service

import (
	"github.com/Viquad/crud-app/internal/domain"
)

type Services struct {
	Account *AccountService
}

func (s *Services) GetAccountService() domain.AccountService {
	return s.Account
}

func NewServices(repo domain.Repositories) *Services {
	return &Services{
		Account: NewAccountService(repo.GetAccountRepository()),
	}
}
