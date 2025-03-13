package repository

import (
	"context"
	"gophermart/internal/repository/mocks"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

type userRow struct {
	id       int
	password string
}

func (m *userRow) Scan(dest ...interface{}) error {
	*(dest[0].(*int)) = m.id
	*(dest[1].(*string)) = m.password
	return nil
}
func TestGetUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mocks.NewMockDBPool(ctrl)

	mockDB.EXPECT().
		QueryRow(gomock.Any(), SelectUser, "test_user").
		Return(&userRow{123, "hashed_password"})

	repo := NewWithPool(mockDB)

	user, err := repo.GetUser(context.Background(), "test_user")
	require.NoError(t, err)
	require.NotNil(t, user)
	require.Equal(t, 123, user.ID)
	require.Equal(t, "hashed_password", user.Password)
}
