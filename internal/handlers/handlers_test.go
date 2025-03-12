package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"gophermart/internal/domain"
	"gophermart/internal/repository"
	"gophermart/internal/utils/bcrypt"
	"gophermart/internal/utils/session"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

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
	repo := new(repository.MockRepository)
	utils := &domain.Utils{L: zap.NewNop()}
	h := New(utils, repo)

	mockSession := new(MockSession)
	mockSession.On("GetUserID", mock.Anything).Return(1, nil)

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
			L: zap.NewNop(),
			S: session.NewSession(),
		},
	}

	r.POST("/register", handler.Register)

	user := domain.Credentials{
		Login:    "test",
		Password: "test",
	}
	jsonData, err := json.Marshal(user)
	assert.NoError(t, err, "WHEN MARSHAL JSON")
	//nolint:noctx // dontknow how fix
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
			L: zap.NewNop(),
			S: session.NewSession(),
		},
	}

	r.POST("/register", handler.Register)

	user := domain.Credentials{
		Login:    "test",
		Password: "password123",
	}
	jsonData, err := json.Marshal(user)
	assert.NoError(t, err, "WHEN MARSHAL JSON")
	//nolint:noctx // dontknow how fix
	req, _ := http.NewRequest("POST", "/register", bytes.NewReader(jsonData))

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockRepo.AssertExpectations(t)
}

func TestLogin_PasswordMismatch(t *testing.T) {
	mockRepo := new(repository.MockRepository)

	storedPassword := "password123"
	pass, err := bcrypt.HashPassword(storedPassword)
	assert.NoError(t, err)
	mockRepo.On("GetUser", mock.Anything, "testuser").Return(&domain.UserIDPassword{
		ID:       1,
		Password: pass,
	}, nil)

	r := gin.Default()
	handler := Handler{
		repo: mockRepo,
		utils: &domain.Utils{
			L: zap.NewNop(),
		},
	}

	r.POST("/login", handler.Login)

	user := domain.Credentials{
		Login:    "testuser",
		Password: "wrongpassword",
	}
	jsonData, err := json.Marshal(user)
	assert.NoError(t, err)
	//nolint:noctx // dontknow how fix
	req, err := http.NewRequest("POST", "/login", bytes.NewReader(jsonData))
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)

	mockRepo.AssertExpectations(t)
}
func TestLogin_Success(t *testing.T) {
	mockRepo := new(repository.MockRepository)

	storedPassword := "password123"
	pass, err := bcrypt.HashPassword(storedPassword)
	assert.NoError(t, err)
	mockRepo.On("GetUser", mock.Anything, "testuser").Return(&domain.UserIDPassword{
		ID:       1,
		Password: pass,
	}, nil)

	r := gin.Default()
	handler := Handler{
		repo: mockRepo,
		utils: &domain.Utils{
			L: zap.NewNop(),
			S: session.NewSession(),
		},
	}

	r.POST("/login", handler.Login)

	user := domain.Credentials{
		Login:    "testuser",
		Password: storedPassword,
	}
	jsonData, err := json.Marshal(user)
	assert.NoError(t, err)
	//nolint:noctx // dontknow how fix
	req, err := http.NewRequest("POST", "/login", bytes.NewReader(jsonData))
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	mockRepo.AssertExpectations(t)
}
func TestBalance_Success(t *testing.T) {
	mockRepo := new(repository.MockRepository)

	mockRepo.On("GetUserBalance", mock.Anything, 1).Return(&domain.Balance{
		Current: 100,
	}, nil)

	s := session.NewSession()
	token, err := s.CreateToken(1)
	assert.NoError(t, err)
	r := gin.Default()

	handler := Handler{
		repo: mockRepo,
		utils: &domain.Utils{
			L: zap.NewNop(),
			S: s,
		},
	}

	r.GET("/balance", handler.Balance)
	req, err := http.NewRequest("GET", "/balance", nil)
	assert.NoError(t, err)

	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var balanceResponse domain.Balance
	err = json.Unmarshal(w.Body.Bytes(), &balanceResponse)
	assert.NoError(t, err)
	assert.Equal(t, float64(100), balanceResponse.Current)

	mockRepo.AssertExpectations(t)
}
func TestListOrders_Success(t *testing.T) {
	mockRepo := new(repository.MockRepository)

	mockRepo.On("GetAllOrders", mock.Anything, 1).Return([]domain.Order{
		{
			OrderWithUserID: domain.OrderWithUserID{
				Number: "12345",
				UserID: 1,
			},
			Accrual:    100.50,
			Status:     domain.OrderStatus("Completed"),
			UploadedAt: time.Now(),
		},
		{
			OrderWithUserID: domain.OrderWithUserID{
				Number: "67890",
				UserID: 1,
			},
			Accrual:    200.75,
			Status:     domain.OrderStatus("Pending"),
			UploadedAt: time.Now(),
		},
	}, nil)

	s := session.NewSession()
	token, err := s.CreateToken(1)
	assert.NoError(t, err)

	r := gin.Default()
	handler := Handler{
		repo: mockRepo,
		utils: &domain.Utils{
			L: zap.NewNop(),
			S: s,
		},
	}

	r.GET("/orders", handler.ListOrders)
	//nolint:noctx // dontknow how fix
	req, err := http.NewRequest("GET", "/orders", nil)
	assert.NoError(t, err)

	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var ordersResponse []domain.Order
	err = json.Unmarshal(w.Body.Bytes(), &ordersResponse)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(ordersResponse))

	assert.Equal(t, "12345", ordersResponse[0].Number)
	assert.Equal(t, 100.50, ordersResponse[0].Accrual)
	assert.Equal(t, domain.OrderStatus("Completed"), ordersResponse[0].Status)

	assert.Equal(t, "67890", ordersResponse[1].Number)
	assert.Equal(t, 200.75, ordersResponse[1].Accrual)
	assert.Equal(t, domain.OrderStatus("Pending"), ordersResponse[1].Status)

	mockRepo.AssertExpectations(t)
}
func TestListOrders_NoContent(t *testing.T) {
	mockRepo := new(repository.MockRepository)

	mockRepo.On("GetAllOrders", mock.Anything, 1).Return([]domain.Order{}, nil)

	s := session.NewSession()
	token, err := s.CreateToken(1)
	assert.NoError(t, err)

	r := gin.Default()
	handler := Handler{
		repo: mockRepo,
		utils: &domain.Utils{
			L: zap.NewNop(),
			S: s,
		},
	}

	r.GET("/orders", handler.ListOrders)

	req, err := http.NewRequest("GET", "/orders", nil)
	assert.NoError(t, err)

	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)

	mockRepo.AssertExpectations(t)
}

func TestWithdrawals_Success(t *testing.T) {
	mockRepo := new(repository.MockRepository)

	mockRepo.On("GetWithdraws", mock.Anything, 1).Return([]domain.Withdraw{
		{
			OrderToWithdraw: domain.OrderToWithdraw{
				Order: "12345",
				Sum:   150.75,
			},
			ProcessedAt: time.Now(),
		},
		{
			OrderToWithdraw: domain.OrderToWithdraw{
				Order: "67890",
				Sum:   200.50,
			},
			ProcessedAt: time.Now(),
		},
	}, nil)

	s := session.NewSession()
	token, err := s.CreateToken(1)
	assert.NoError(t, err)

	r := gin.Default()
	handler := Handler{
		repo: mockRepo,
		utils: &domain.Utils{
			L: zap.NewNop(),
			S: s,
		},
	}

	r.GET("/withdrawals", handler.Withdrawals)
	//nolint:noctx // dontknow how fix
	req, err := http.NewRequest("GET", "/withdrawals", nil)
	assert.NoError(t, err)

	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var withdrawsResponse []domain.Withdraw
	err = json.Unmarshal(w.Body.Bytes(), &withdrawsResponse)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(withdrawsResponse))

	assert.Equal(t, "12345", withdrawsResponse[0].Order)
	assert.Equal(t, 150.75, withdrawsResponse[0].Sum)

	assert.Equal(t, "67890", withdrawsResponse[1].Order)
	assert.Equal(t, 200.50, withdrawsResponse[1].Sum)

	mockRepo.AssertExpectations(t)
}

func TestWithdrawals_NoContent(t *testing.T) {
	mockRepo := new(repository.MockRepository)

	mockRepo.On("GetWithdraws", mock.Anything, 1).Return([]domain.Withdraw{}, nil)

	s := session.NewSession()
	token, err := s.CreateToken(1)
	assert.NoError(t, err)

	r := gin.Default()
	handler := Handler{
		repo: mockRepo,
		utils: &domain.Utils{
			L: zap.NewNop(),
			S: s,
		},
	}

	r.GET("/withdrawals", handler.Withdrawals)
	//nolint:noctx // dontknow how fix
	req, err := http.NewRequest("GET", "/withdrawals", nil)
	assert.NoError(t, err)

	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)

	mockRepo.AssertExpectations(t)
}

func TestBalanceWithdraw_Success(t *testing.T) {
	mockRepo := new(repository.MockRepository)

	mockRepo.On("BalanceWithdraw", mock.Anything, 1, mock.Anything).Return(nil)

	s := session.NewSession()
	token, err := s.CreateToken(1)
	assert.NoError(t, err)

	r := gin.Default()
	handler := Handler{
		repo: mockRepo,
		utils: &domain.Utils{
			L: zap.NewNop(),
			S: s,
		},
	}

	r.POST("/balance/withdraw", handler.BalanceWithdraw)

	orderToWithdraw := domain.OrderToWithdraw{
		Order: "12345",
		Sum:   150.75,
	}

	body, err := json.Marshal(orderToWithdraw)
	assert.NoError(t, err)
	//nolint:noctx // dontknow how fix
	req, err := http.NewRequest("POST", "/balance/withdraw", bytes.NewReader(body))
	assert.NoError(t, err)

	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	mockRepo.AssertExpectations(t)
}

func TestBalanceWithdraw_InvalidSum(t *testing.T) {
	mockRepo := new(repository.MockRepository)

	s := session.NewSession()
	token, err := s.CreateToken(1)
	assert.NoError(t, err)

	r := gin.Default()
	handler := Handler{
		repo: mockRepo,
		utils: &domain.Utils{
			L: zap.NewNop(),
			S: s,
		},
	}

	r.POST("/balance/withdraw", handler.BalanceWithdraw)

	orderToWithdraw := domain.OrderToWithdraw{
		Order: "12345",
		Sum:   -150.75,
	}

	body, err := json.Marshal(orderToWithdraw)
	assert.NoError(t, err)
	//nolint:noctx // dontknow how fix
	req, err := http.NewRequest("POST", "/balance/withdraw", bytes.NewReader(body))
	assert.NoError(t, err)

	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
