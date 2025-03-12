package handlers

import (
	"context"
	"gophermart/internal/domain"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

// MockRepository is a mock implementation of the Repository interface
type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) SaveUser(ctx context.Context, user *domain.Credentials) (int, error) {
	args := m.Called(ctx, user)
	return args.Int(0), args.Error(1)
}

func (m *MockRepository) GetUser(ctx context.Context, login string) (*domain.UserIDPassword, error) {
	args := m.Called(ctx, login)
	return args.Get(0).(*domain.UserIDPassword), args.Error(1)
}

func (m *MockRepository) GetOrderWithUserID(ctx context.Context, number string) (*domain.OrderWithUserID, error) {
	args := m.Called(ctx, number)
	return args.Get(0).(*domain.OrderWithUserID), args.Error(1)
}

func (m *MockRepository) CreateOrder(ctx context.Context, data *domain.OrderWithUserID) error {
	args := m.Called(ctx, data)
	return args.Error(0)
}

func (m *MockRepository) GetAllOrders(ctx context.Context, userID int) ([]domain.Order, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]domain.Order), args.Error(1)
}

func (m *MockRepository) GetUserBalance(ctx context.Context, userID int) (*domain.Balance, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(*domain.Balance), args.Error(1)
}

func (m *MockRepository) BalanceWithdraw(ctx context.Context, userID int, withdraw *domain.OrderToWithdraw) error {
	args := m.Called(ctx, userID, withdraw)
	return args.Error(0)
}

func (m *MockRepository) GetWithdraws(ctx context.Context, userID int) ([]domain.Withdraw, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]domain.Withdraw), args.Error(1)
}
func TestNew(t *testing.T) {
	utils := &domain.Utils{}
	repo := new(MockRepository)

	h := New(utils, repo)

	assert.NotNil(t, h)
	assert.Equal(t, utils, h.utils)
	assert.Equal(t, repo, h.repo)
}

type MockSession struct {
	mock.Mock
}

func (m *MockSession) GetUserID(authHeader string) (int, error) {
	args := m.Called(authHeader)
	return args.Int(0), args.Error(1)
}
func TestListOrders(t *testing.T) {
	// Create a mock repository
	repo := new(MockRepository)
	utils := &domain.Utils{L: zap.NewNop()} // Replace with your logger
	h := New(utils, repo)

	// Create a mock session
	mockSession := new(MockSession)
	mockSession.On("GetUserID", mock.Anything).Return(1, nil)

	// Create a gin context
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/orders", func(c *gin.Context) {
		requestUserID, _ := mockSession.GetUserID(c.Request.Header.Get("Authorization"))
		allOrders, err := h.repo.GetAllOrders(context.TODO(), requestUserID)
		if err != nil {
			h.utils.L.Warn("error getting allOrders", zap.Error(err))
			c.Status(http.StatusInternalServerError)
			return
		}
		if len(allOrders) == 0 {
			c.Status(http.StatusNoContent)
			return
		}
		c.JSON(http.StatusOK, allOrders)
	})

	// Test case: successful retrieval of orders
	t.Run("successful retrieval of orders", func(t *testing.T) {
		orders := []domain.Order{
			{OrderWithUserID: domain.OrderWithUserID{Number: "1", UserID: 1}},
			{OrderWithUserID: domain.OrderWithUserID{Number: "2", UserID: 1}},
		}
		repo.On("GetAllOrders", mock.Anything, 1).Return(orders, nil)

		req, _ := http.NewRequest(http.MethodGet, "/orders", nil)
		req.Header.Set("Authorization", "Bearer token")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, `[
  {
    "number": "1",
    "status": "",
    "uploaded_at": "0001-01-01T00:00:00Z"
  },
  {
    "number": "2",
    "status": "",
    "uploaded_at": "0001-01-01T00:00:00Z"
  }
]`, w.Body.String())
		repo.AssertExpectations(t)
	})

	// Test case: no orders found
	t.Run("no orders found", func(t *testing.T) {
		repo.On("GetAllOrders", mock.Anything, 1).Return([]domain.Order{}, nil)

		req, _ := http.NewRequest(http.MethodGet, "/orders", nil)
		req.Header.Set("Authorization", "Bearer token")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNoContent, w.Code)
		repo.AssertExpectations(t)
	})
}
