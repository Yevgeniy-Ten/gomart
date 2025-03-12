package repository

import (
	"context"
	"gophermart/internal/domain"
	"gophermart/internal/repository/mocks"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

// 🔹 Фейковая структура, которая реализует `pgx.Row`
type balanceRow struct {
	current  float64
	withdraw float64
	err      error
}

// Реализуем метод Scan(), который заполняет переданные аргументы
func (m *balanceRow) Scan(dest ...interface{}) error {
	if m.err != nil {
		return m.err
	}
	*(dest[0].(*float64)) = m.current
	*(dest[1].(*float64)) = m.withdraw
	return nil
}

func TestGetUserBalance(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mocks.NewMockDBPool(ctrl)

	// 🔹 Данные для теста
	userID := 1
	expectedBalance := domain.Balance{
		Current:  500.75,
		Withdraw: 100.25,
	}

	// 🔹 Мокаем QueryRow(), чтобы он вернул кастомный `balanceRow`
	mockDB.EXPECT().
		QueryRow(gomock.Any(), SelectUserBalance, userID).
		Return(&balanceRow{
			current:  expectedBalance.Current,
			withdraw: expectedBalance.Withdraw,
		})

	repo := NewWithPool(mockDB)
	balance, err := repo.GetUserBalance(context.Background(), userID)

	// 🔹 Проверяем, что нет ошибки и баланс совпадает с ожидаемым
	require.NoError(t, err)
	require.Equal(t, expectedBalance, *balance)
}
