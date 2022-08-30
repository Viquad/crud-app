package psql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/Viquad/crud-app/internal/domain"
)

type AccountRepository struct {
	db *sql.DB
}

func NewAccountRepository(db *sql.DB) *AccountRepository {
	return &AccountRepository{
		db: db,
	}
}

func (b *AccountRepository) Create(ctx context.Context, inp domain.AccountCreateInput) (*domain.Account, error) {
	userId, ok := ctx.Value(domain.UserIdKey).(int64)
	if !ok {
		return nil, domain.ErrInvalidId
	}

	account := domain.Account{
		UserId:   userId,
		Balance:  inp.Balance,
		Currency: inp.Currency,
	}

	query := "INSERT INTO accounts (user_id, balance, currency) VALUES ($1, $2, $3) RETURNING id, last_update"
	err := b.db.QueryRowContext(ctx, query, userId, inp.Balance, inp.Currency).
		Scan(&account.Id, &account.LastUpdate)

	if err != nil {
		return nil, err
	}

	return &account, err
}

func (b *AccountRepository) GetById(ctx context.Context, id int64) (*domain.Account, error) {
	var account domain.Account

	userId, ok := ctx.Value(domain.UserIdKey).(int64)
	if !ok {
		return nil, domain.ErrInvalidId
	}

	query := "SELECT id, user_id, balance, currency, last_update FROM accounts WHERE id = $1 AND user_id = $2"
	row := b.db.QueryRowContext(ctx, query, id, userId)

	if err := row.Scan(&account.Id, &account.UserId, &account.Balance, &account.Currency, &account.LastUpdate); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNotExist
		}
		return nil, err
	}

	return &account, nil
}

func (b *AccountRepository) List(ctx context.Context) ([]domain.Account, error) {
	var accounts []domain.Account

	userId, ok := ctx.Value(domain.UserIdKey).(int64)
	if !ok {
		return nil, domain.ErrInvalidId
	}

	query := "SELECT id, user_id, balance, currency, last_update FROM accounts WHERE user_id = $1"
	rows, err := b.db.QueryContext(ctx, query, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var account domain.Account
		if err := rows.Scan(&account.Id, &account.UserId, &account.Balance, &account.Currency, &account.LastUpdate); err != nil {
			return nil, err
		}

		accounts = append(accounts, account)
	}

	return accounts, nil
}

func (b *AccountRepository) UpdateById(ctx context.Context, id int64, inp domain.AccountUpdateInput) (*domain.Account, error) {
	var (
		account   domain.Account
		setValues []string
		args      []interface{}
		argIndex  = 1
		addArg    = func(i interface{}, arg string) {
			setValues = append(setValues, fmt.Sprintf("%s=$%d", arg, argIndex))
			args = append(args, i)
			argIndex++
		}
	)

	userId, ok := ctx.Value(domain.UserIdKey).(int64)
	if !ok {
		return nil, domain.ErrInvalidId
	}

	if inp.Balance != nil {
		addArg(*inp.Balance, "balance")
	}
	addArg("now()", "last_update")

	setQuery := strings.Join(setValues, ", ")
	query := fmt.Sprintf("UPDATE accounts SET %s WHERE id=$%d AND user_id=$%d RETURNING id, user_id, balance, currency, last_update", setQuery, argIndex, argIndex+1)
	argIndex++
	args = append(args, id, userId)

	row := b.db.QueryRowContext(ctx, query, args...)
	err := row.Scan(&account.Id, &account.UserId, &account.Balance, &account.Currency, &account.LastUpdate)
	if err != nil {
		return nil, domain.ErrUpdateFailed
	}

	return &account, nil
}

func (b *AccountRepository) DeleteById(ctx context.Context, id int64) error {
	userId, ok := ctx.Value(domain.UserIdKey).(int64)
	if !ok {
		return domain.ErrInvalidId
	}

	res, err := b.db.ExecContext(ctx, "DELETE FROM accounts WHERE id=$1 AND user_id=$2", id, userId)
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
