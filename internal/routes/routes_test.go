package routes

import (
	"gophermart/internal/utils"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHasUserID(t *testing.T) {
	baseURI := "/api/user/"
	needUserID := []struct {
		url    string
		method string
	}{
		{baseURI + "orders", "POST"},
		{baseURI + "orders", "GET"},
		{baseURI + "balance", "GET"},
		{baseURI + "balance/withdraw", "POST"},
		{baseURI + "withdrawals", "GET"},
	}
	u, err := utils.New(nil)
	assert.NoError(t, err, "WHEN INIT UTILS")
	r := Init(u, nil)
	for _, tt := range needUserID {
		t.Run(tt.url, func(t *testing.T) {
			request := httptest.NewRequest(tt.method, tt.url, nil)
			recorder := httptest.NewRecorder()
			r.ServeHTTP(recorder, request)
			result := recorder.Result()
			defer result.Body.Close()
			assert.Equal(t, http.StatusUnauthorized, result.StatusCode)
		})
	}
}
