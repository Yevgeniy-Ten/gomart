package repository

import (
	"context"
	"gophermart/internal/domain"
	"gophermart/internal/repository/mocks"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type orderRow struct {
	number string
	userID int
}

func (m *orderRow) Scan(dest ...interface{}) error {
	*(dest[0].(*string)) = m.number
	*(dest[1].(*int)) = m.userID
	return nil
}

func TestGetOrderWithUserID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mocks.NewMockDBPool(ctrl)

	expectedOrder := domain.OrderWithUserID{
		Number: "123456",
		UserID: 42,
	}

	mockDB.EXPECT().
		QueryRow(gomock.Any(), SelectOrder, "123456").
		Return(&orderRow{expectedOrder.Number, expectedOrder.UserID}) // Используем кастомный orderRow

	repo := NewWithPool(mockDB)

	order, err := repo.GetOrderWithUserID(context.Background(), "123456")
	require.NoError(t, err)
	require.NotNil(t, order)
	require.Equal(t, expectedOrder.Number, order.Number)
	require.Equal(t, expectedOrder.UserID, order.UserID)
}

func TestCreateOrder(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mocks.NewMockDBPool(ctrl)

	order := &domain.OrderWithUserID{
		Number: "123456",
		UserID: 42,
	}

	mockDB.EXPECT().
		Exec(gomock.Any(), InsertOrder, gomock.Any(), gomock.Any()).
		Return(pgconn.NewCommandTag("INSERT 1"), nil)

	repo := NewWithPool(mockDB)

	err := repo.CreateOrder(context.Background(), order)
	require.NoError(t, err)
}

func TestGroupAccrualsByUserID(t *testing.T) {
	accrual1 := 10.5
	accrual2 := 20.0
	accrual3 := 15.5

	orders := []*domain.OrderWithAccrual{
		{AccrualResponse: domain.AccrualResponse{Accrual: &accrual1}, OrderWithUserID: domain.OrderWithUserID{UserID: 1}},
		{AccrualResponse: domain.AccrualResponse{Accrual: &accrual2}, OrderWithUserID: domain.OrderWithUserID{UserID: 2}},
		{AccrualResponse: domain.AccrualResponse{Accrual: &accrual3}, OrderWithUserID: domain.OrderWithUserID{UserID: 1}},
		{AccrualResponse: domain.AccrualResponse{Accrual: nil}, OrderWithUserID: domain.OrderWithUserID{UserID: 3}},
	}

	expected := map[int]float64{
		1: 26.0, // 10.5 + 15.5
		2: 20.0,
	}

	result := GroupAccrualsByUserID(orders)
	assert.Equal(t, expected, result, "Accruals were not grouped correctly by user ID")
}
func TestGetNewOrders(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mocks.NewMockDBPool(ctrl)
	mockRows := mocks.NewMockRows(ctrl)

	// Упорядоченный вызов Next(), Scan(), Close()
	gomock.InOrder(
		mockRows.EXPECT().Next().Return(true),
		mockRows.EXPECT().Scan(gomock.Any(), gomock.Any()).Return(nil),
		mockRows.EXPECT().Next().Return(true),
		mockRows.EXPECT().Scan(gomock.Any(), gomock.Any()).Return(nil),
		mockRows.EXPECT().Next().Return(false), // Конец данных
		mockRows.EXPECT().Close(),              // Close() не возвращает значения
	)

	mockDB.EXPECT().
		Query(gomock.Any(), SelectOrdersByNewAndProcessingStatus, gomock.Any()).
		Return(mockRows, nil)

	repo := Repo{conn: mockDB}
	orders, err := repo.GetNewOrders(context.Background(), 2)

	require.NoError(t, err)
	require.Len(t, orders, 2)
}
