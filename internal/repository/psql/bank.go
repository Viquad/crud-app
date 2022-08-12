package psql

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/Viquad/crud-app/internal/domain"
)

type Bank struct {
	db *sql.DB
}

func NewBank(db *sql.DB) *Bank {
	return &Bank{
		db: db,
	}
}

func (b *Bank) CreateAccount(ctx context.Context, account domain.Account) error {
	_, err := b.db.Exec("INSERT INTO accounts (first_name, last_name, balance, currency) VALUES ($1, $2, $3, $4)",
		account.FirstName, account.LastName, account.Balance, account.Currency)

	return err
}

func (b *Bank) ReadAccount(ctx context.Context, id int64) (domain.Account, error) {
	var account domain.Account
	err := b.db.QueryRow("SELECT id, first_name, last_name, balance, currency, last_update FROM accounts WHERE id = $1", id).
		Scan(&account.Id, &account.FirstName, &account.LastName, &account.Balance, &account.Currency, &account.LastUpdate)

	return account, err
}

func (b *Bank) ReadAllAccounts(ctx context.Context) ([]domain.Account, error) {
	var accounts []domain.Account
	rows, err := b.db.Query("SELECT id, first_name, last_name, balance, currency, last_update FROM accounts")
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var account domain.Account
		if err := rows.Scan(&account.Id, &account.FirstName, &account.LastName, &account.Balance, &account.Currency, &account.LastUpdate); err != nil {
			return nil, err
		}

		accounts = append(accounts, account)
	}

	return accounts, rows.Err()
}

func (b *Bank) UpdateAccount(ctx context.Context, id int64, inp domain.AccountUpdateInput) error {
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argIndex := 1

	addArg := func(i interface{}, arg string) {
		setValues = append(setValues, fmt.Sprintf("%s=$%d", arg, argIndex))
		args = append(args, i)
		argIndex++
	}

	if inp.FirstName != nil {
		addArg(*inp.FirstName, "first_name")
	}
	if inp.LastName != nil {
		addArg(*inp.LastName, "last_name")
	}
	if inp.Balance != nil {
		addArg(*inp.Balance, "balance")
	}
	if inp.Currency != nil {
		addArg(*inp.Currency, "currency")
	}
	addArg(time.Now().Format(time.RFC3339), "last_update")

	setQuery := strings.Join(setValues, ", ")

	query := fmt.Sprintf("UPDATE accounts SET %s WHERE id=$%d", setQuery, argIndex)
	args = append(args, id)

	_, err := b.db.Exec(query, args...)
	return err
}

func (b *Bank) DeleteAccount(ctx context.Context, id int64) error {
	_, err := b.db.Exec("DELETE FROM accounts WHERE id=$1", id)

	return err
}
