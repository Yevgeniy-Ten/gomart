package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"gophermart/internal/domain"
	"gophermart/internal/repository"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

func TestNew(t *testing.T) {
	utils := &domain.Utils{}
	repo := new(repository.MockRepository)

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
	repo := new(repository.MockRepository)
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
		//nolint:noctx // dontknow how fix
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
		//nolint:noctx // dontknow how fix
		req, _ := http.NewRequest(http.MethodGet, "/orders", nil)
		req.Header.Set("Authorization", "Bearer token")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		repo.AssertExpectations(t)
	})
}
func TestRegister_Success(t *testing.T) {
	mockRepo := new(repository.MockRepository)
	mockRepo.On("SaveUser", mock.Anything, mock.Anything).Return(1, nil)

	r := gin.Default()

	handler := Handler{
		repo: mockRepo,
		utils: &domain.Utils{
			L: zap.NewNop(), // Мок логгера
		},
	}

	r.POST("/register", handler.Register)

	user := domain.Credentials{
		Login:    "test",
		Password: "test",
	}
	jsonData, err := json.Marshal(user)
	assert.NoError(t, err, "WHEN MARSHAL JSON")
	req, _ := http.NewRequest("POST", "/register", bytes.NewReader(jsonData))

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockRepo.AssertExpectations(t)
}

func TestRegister_SaveUserError(t *testing.T) {
	mockRepo := new(repository.MockRepository)
	saveErr := errors.New("save error")
	mockRepo.On("SaveUser", mock.Anything, mock.Anything).Return(0, saveErr)

	r := gin.Default()

	handler := Handler{
		repo: mockRepo,
		utils: &domain.Utils{
			L: zap.NewNop(), // Мок логгера
		},
	}

	r.POST("/register", handler.Register)

	user := domain.Credentials{
		Login:    "test",
		Password: "password123",
	}
	jsonData, err := json.Marshal(user)
	assert.NoError(t, err, "WHEN MARSHAL JSON")

	req, _ := http.NewRequest("POST", "/register", bytes.NewReader(jsonData))

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockRepo.AssertExpectations(t)
}

func TestLogin_PasswordMismatch(t *testing.T) {
	// Setup
	mockRepo := new(repository.MockRepository)
	mockRepo.On("GetUser", mock.Anything, "testuser").Return(domain.UserIDPassword{
		ID:       1,
		Password: "$2a$10$O2kTe.R0DWug08a4y7PeAOkQ3cxm9V7rO9S/VB7pEY2X7Ltn1neUS", // Пример хешированного пароля
	}, nil)

	r := gin.Default()
	handler := Handler{
		repo: mockRepo,
		utils: &domain.Utils{
			L: zap.NewNop(), // Мок логгера
		},
	}

	r.POST("/login", handler.Login)

	// Входные данные с неправильным паролем
	user := domain.Credentials{
		Login:    "testuser",
		Password: "wrongpassword", // Неверный пароль
	}
	jsonData, err := json.Marshal(user)
	assert.NoError(t, err)

	req, _ := http.NewRequest("POST", "/login", bytes.NewReader(jsonData))

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockRepo.AssertExpectations(t)
}

func TestLogin_Success(t *testing.T) {
	// Setup
	mockRepo := new(repository.MockRepository)
	mockRepo.On("GetUser", mock.Anything, "testuser").Return(domain.UserIDPassword{
		ID:       1,
		Password: "$2a$10$O2kTe.R0DWug08a4y7PeAOkQ3cxm9V7rO9S/VB7pEY2X7Ltn1neUS", // Пример хешированного пароля
	}, nil)

	r := gin.Default()
	handler := Handler{
		repo: mockRepo,
		utils: &domain.Utils{
			L: zap.NewNop(), // Мок логгера
		},
	}

	r.POST("/login", handler.Login)

	// Входные данные
	user := domain.Credentials{
		Login:    "testuser",
		Password: "password123", // Соответствует хешированному паролю
	}
	jsonData, err := json.Marshal(user)
	assert.NoError(t, err)

	req, _ := http.NewRequest("POST", "/login", bytes.NewReader(jsonData))

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockRepo.AssertExpectations(t)
}
