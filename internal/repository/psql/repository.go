package psql

import (
	"database/sql"

	"github.com/Viquad/crud-app/internal/domain"
)

type Repositories struct {
	Account *AccountRepository
}

func (r *Repositories) GetAccountRepository() domain.AccountRepository {
	return r.Account
}

func NewRepositories(db *sql.DB) *Repositories {
	return &Repositories{
		Account: NewAccount(db),
	}
}
