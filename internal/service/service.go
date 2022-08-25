package service

import (
	"time"

	"github.com/Viquad/crud-app/internal/domain"
	cache "github.com/Viquad/simple-cache"
)

type Repositories interface {
	GetAccountRepository() domain.AccountRepository
	GetUserRepository() domain.UserRepository
}

type PasswordHasher interface {
	Hash(password string) (string, error)
}

type Services struct {
	accountService *AccountService
	userService    *UserService
}

func (ss *Services) GetAccountService() domain.AccountService {
	return ss.accountService
}

func (ss *Services) GetUserService() domain.UserService {
	return ss.userService
}

func NewServices(repo Repositories, cache cache.Cache, cachettl time.Duration, hasher PasswordHasher, secret []byte, tokenttl time.Duration) *Services {
	return &Services{
		accountService: NewAccountService(repo.GetAccountRepository(), cache, cachettl),
		userService:    NewUserService(repo.GetUserRepository(), hasher, secret, tokenttl),
	}
}
