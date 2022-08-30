package psql

import (
	"database/sql"

	"github.com/Viquad/crud-app/internal/domain"
)

type Repositories struct {
	accountRepository *AccountRepository
	userRepository    *UserRepository
	tokenRepository   *TokenRepository
}

func (rs *Repositories) GetAccountRepository() domain.AccountRepository {
	return rs.accountRepository
}

func (rs *Repositories) GetUserRepository() domain.UserRepository {
	return rs.userRepository
}

func (rs *Repositories) GetTokenRepository() domain.TokenRepository {
	return rs.tokenRepository
}

func NewRepositories(db *sql.DB) *Repositories {
	return &Repositories{
		accountRepository: NewAccountRepository(db),
		userRepository:    NewUserRepository(db),
		tokenRepository:   NewTokenRepository(db),
	}
}
