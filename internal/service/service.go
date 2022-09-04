package service

import (
	"time"

	"github.com/Viquad/crud-app/internal/domain"
	"github.com/Viquad/crud-audit-service/pkg/domain/audit"
	cache "github.com/Viquad/simple-cache"
	"github.com/gorilla/sessions"
)

type Repositories interface {
	GetAccountRepository() domain.AccountRepository
	GetUserRepository() domain.UserRepository
	GetTokenRepository() domain.TokenRepository
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

func NewServices(repo Repositories, cache cache.Cache, store sessions.Store, hasher PasswordHasher, audit audit.AuditServiceClient, secret []byte, cachettl, accessttl, refreshttl time.Duration) *Services {
	return &Services{
		accountService: NewAccountService(repo, cache, cachettl),
		userService:    NewUserService(repo, store, audit, hasher, secret, accessttl, refreshttl),
	}
}
