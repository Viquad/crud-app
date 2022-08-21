package service

import (
	"time"

	"github.com/Viquad/crud-app/internal/domain"
	cache "github.com/Viquad/simple-cache"
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

func NewServices(repo Repositories, cache cache.Cache, ttl time.Duration) *Services {
	return &Services{
		accountService: NewAccountService(repo.GetAccountRepository(), cache, ttl),
	}
}
