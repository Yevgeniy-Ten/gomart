package repository

import (
	"context"
	"errors"
	"gophermart/internal/domain"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
)

const (
	CreateBalance     = "INSERT INTO balance (user_id) VALUES ($1)"
	SelectUserBalance = "SELECT current,withdrawn FROM balance WHERE user_id = $1"
)

func (d *Repo) GetUserBalance(ctx context.Context, userID int) (*domain.Balance, error) {
	var balance domain.Balance
	err := d.conn.QueryRow(ctx, SelectUserBalance, userID).Scan(&balance.Current, &balance.Withdraw)
	return &balance, err
}

const InsertIntoPayments = `INSERT INTO payments (sum, user_id, "order") VALUES ($1, $2, $3)`
const UpdateBalance = `UPDATE balance SET current = current - $1, withdrawn = withdrawn + $1 WHERE user_id = $2`

func (d *Repo) BalanceWithdraw(ctx context.Context, userID int, withdraw *domain.OrderToWithdraw) error {
	tx, err := d.conn.Begin(ctx)
	if err != nil {
		return err
	}
	//nolint:errcheck // ignore error because it's not important
	defer tx.Rollback(ctx)
	_, err = tx.Exec(ctx, UpdateBalance, withdraw.Sum, userID)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.CheckViolation {
			return NewShouldBePositiveError()
		}
		return errors.New("error insert into payments" + err.Error())
	}
	_, err = tx.Exec(ctx, InsertIntoPayments, withdraw.Sum, userID, withdraw.Order)
	if err != nil {
		return errors.New("error update balance" + err.Error())
	}

	if err = tx.Commit(ctx); err != nil {
		return errors.New("error commit" + err.Error())
	}
	return nil
}

const SelectWithdraws = `SELECT sum,"order",processed_at FROM payments WHERE user_id = $1`

func (d *Repo) GetWithdraws(ctx context.Context, userID int) ([]domain.Withdraw, error) {
	rows, err := d.conn.Query(ctx, SelectWithdraws, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var withdraws []domain.Withdraw
	for rows.Next() {
		var withdraw domain.Withdraw
		if err := rows.Scan(&withdraw.Sum, &withdraw.Order, &withdraw.ProcessedAt); err != nil {
			return nil, err
		}
		withdraws = append(withdraws, withdraw)
	}
	return withdraws, nil
}
