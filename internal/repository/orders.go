package repository

import (
	"context"
	"errors"
	"fmt"
	"gophermart/internal/domain"

	"github.com/jackc/pgx/v5"
)

const SelectOrder = "SELECT number, user_id FROM orders WHERE number = $1"

func (d *Repo) GetOrderWithUserID(ctx context.Context, number string) (*domain.OrderWithUserID, error) {
	var order domain.OrderWithUserID
	err := d.conn.QueryRow(ctx, SelectOrder, number).Scan(&order.Number, &order.UserID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, NewNotFoundError(number)
		}
		return nil, err
	}
	return &order, nil
}

const InsertOrder = "INSERT INTO orders (number, user_id) VALUES ($1, $2)"

func (d *Repo) CreateOrder(ctx context.Context, data *domain.OrderWithUserID) error {
	_, err := d.conn.Exec(ctx, InsertOrder, data.Number, data.UserID)
	return err
}

const SelectAllOrders = "SELECT number,accrual,status,uploaded_at FROM orders WHERE user_id = $1"

func (d *Repo) GetAllOrders(ctx context.Context, userID int) ([]domain.Order, error) {
	rows, err := d.conn.Query(ctx, SelectAllOrders, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []domain.Order
	for rows.Next() {
		var order domain.Order
		var accrual *float64
		if err := rows.Scan(&order.Number, &accrual, &order.Status, &order.UploadedAt); err != nil {
			return nil, err
		}
		if accrual != nil {
			order.Accrual = *accrual
		}
		orders = append(orders, order)
	}
	return orders, nil
}

const SelectOrdersByNewAndProcessingStatus = `SELECT number, user_id FROM orders 
                       WHERE status in ('NEW','PROCESSING') ORDER BY uploaded_at LIMIT $1`

func (d *Repo) GetNewOrders(ctx context.Context, limit int) ([]*domain.OrderWithUserID, error) {
	rows, err := d.conn.Query(ctx, SelectOrdersByNewAndProcessingStatus, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var orders []*domain.OrderWithUserID
	for rows.Next() {
		var number string
		var userID int
		if err := rows.Scan(&number, &userID); err != nil {
			return nil, err
		}
		orders = append(orders, &domain.OrderWithUserID{Number: number, UserID: userID})
	}
	return orders, nil
}

func GroupAccrualsByUserID(accruals []*domain.OrderWithAccrual) map[int]float64 {
	grouped := make(map[int]float64)
	for _, a := range accruals {
		if a.Accrual == nil {
			continue
		}
		grouped[a.UserID] += *a.Accrual
	}
	return grouped
}

const UpdateOrderStatus = "UPDATE orders SET status = $1, accrual = $2 WHERE number = $3"

func (d *Repo) UpdateOrdersWithAccrual(ctx context.Context, accruals []*domain.OrderWithAccrual) error {
	batch := &pgx.Batch{}
	for _, a := range accruals {
		batch.Queue(UpdateOrderStatus, a.Status, a.Accrual, a.Number)
	}

	userAccruals := GroupAccrualsByUserID(accruals)
	for userID, accrual := range userAccruals {
		if accrual == 0 {
			continue
		}
		batch.Queue(AddToBalance, accrual, userID)
	}
	br := d.conn.SendBatch(ctx, batch)
	defer br.Close()
	for range accruals {
		if _, err := br.Exec(); err != nil {
			return fmt.Errorf("failed to execute batch update: %w", err)
		}
	}
	return nil
}
