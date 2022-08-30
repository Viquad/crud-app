package psql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/Viquad/crud-app/internal/domain"
)

type AccountRepository struct {
	db *sql.DB
}

func NewAccount(db *sql.DB) *AccountRepository {
	return &AccountRepository{
		db: db,
	}
}

func (b *AccountRepository) Create(ctx context.Context, account domain.Account) (*domain.Account, error) {
	var id int64
	var update time.Time
	query := "INSERT INTO accounts (first_name, last_name, balance, currency) VALUES ($1, $2, $3, $4) RETURNING id, last_update"
	err := b.db.QueryRowContext(ctx, query, account.FirstName, account.LastName, account.Balance, account.Currency).
		Scan(&id, &update)

	if err != nil {
		return nil, err
	}

	account.Id = id
	account.LastUpdate = update

	return &account, err
}

func (b *AccountRepository) GetById(ctx context.Context, id int64) (*domain.Account, error) {
	var account domain.Account
	query := "SELECT id, first_name, last_name, balance, currency, last_update FROM accounts WHERE id = $1"
	row := b.db.QueryRowContext(ctx, query, id)

	if err := row.Scan(&account.Id, &account.FirstName, &account.LastName, &account.Balance, &account.Currency, &account.LastUpdate); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNotExist
		}
		return nil, err
	}

	return &account, nil
}

func (b *AccountRepository) All(ctx context.Context) ([]domain.Account, error) {
	var accounts []domain.Account
	query := "SELECT id, first_name, last_name, balance, currency, last_update FROM accounts"
	rows, err := b.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var account domain.Account
		if err := rows.Scan(&account.Id, &account.FirstName, &account.LastName, &account.Balance, &account.Currency, &account.LastUpdate); err != nil {
			return nil, err
		}

		accounts = append(accounts, account)
	}

	return accounts, nil
}

func (b *AccountRepository) Update(ctx context.Context, id int64, inp domain.AccountUpdateInput) (*domain.Account, error) {
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argIndex := 1
	account := new(domain.Account)

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
	addArg("now()", "last_update")

	setQuery := strings.Join(setValues, ", ")
	query := fmt.Sprintf("UPDATE accounts SET %s WHERE id=$%d RETURNING id, first_name, last_name, balance, currency, last_update", setQuery, argIndex)
	args = append(args, id)

	row := b.db.QueryRowContext(ctx, query, args...)
	err := row.Scan(&account.Id, &account.FirstName, &account.LastName, &account.Balance, &account.Currency, &account.LastUpdate)
	if err != nil {
		return nil, domain.ErrUpdateFailed
	}

	return account, nil
}

func (b *AccountRepository) Delete(ctx context.Context, id int64) error {
	res, err := b.db.ExecContext(ctx, "DELETE FROM accounts WHERE id=$1", id)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return domain.ErrDeleteFailed
	}

	return nil
}
