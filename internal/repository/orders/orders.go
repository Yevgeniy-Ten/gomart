package orders

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"gophermart/internal/domain"
)

type OrderRepo struct {
	conn *pgxpool.Pool
}

func New(conn *pgxpool.Pool) *OrderRepo {
	return &OrderRepo{conn: conn}
}

const SelectOrder = "SELECT number, user_id FROM orders WHERE number = $1"

func (o *OrderRepo) GetOrderWithUserID(ctx context.Context, number string) (*domain.OrderWithUserID, error) {
	var order domain.OrderWithUserID
	err := o.conn.QueryRow(ctx, SelectOrder, number).Scan(&order.Number, &order.UserID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, NewNotFoundError(number)
		}
		return nil, err
	}
	return &order, nil
}

const InsertOrder = "INSERT INTO orders (number, user_id) VALUES ($1, $2)"

func (o *OrderRepo) CreateOrder(ctx context.Context, data *domain.OrderWithUserID) error {
	_, err := o.conn.Exec(ctx, InsertOrder, data.Number, data.UserID)
	return err
}

const SelectAllOrders = "SELECT number,accrual,status,uploaded_at FROM orders WHERE user_id = $1"

func (o *OrderRepo) GetAllOrders(ctx context.Context, userID int) ([]domain.Order, error) {
	rows, err := o.conn.Query(ctx, SelectAllOrders, userID)
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
