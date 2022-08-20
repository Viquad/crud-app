package psql

import (
	"database/sql"

	"github.com/Viquad/crud-app/internal/domain"
)

type Repositories struct {
	accountRepository *AccountRepository
}

func (rs *Repositories) GetAccountRepository() domain.AccountRepository {
	return rs.accountRepository
}

func NewRepositories(db *sql.DB) *Repositories {
	return &Repositories{
		accountRepository: NewAccountRepository(db),
	}
}
