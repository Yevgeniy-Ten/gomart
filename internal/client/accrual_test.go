package client

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"gophermart/internal/domain"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetAccrualOrderStatus(t *testing.T) {
	tests := []struct {
		name          string
		statusCode    int
		retryAfter    string
		expectedError error
	}{
		{
			name:          "Status OK",
			statusCode:    http.StatusOK,
			expectedError: nil,
		},
		{
			name:          "Too Many Requests",
			statusCode:    http.StatusTooManyRequests,
			retryAfter:    "10",
			expectedError: &domain.TooManyRequestsError{RetryAfter: 10},
		},
		{
			name:          "Invalid Retry-After Header",
			statusCode:    http.StatusTooManyRequests,
			retryAfter:    "invalid",
			expectedError: errors.New("invalid retry-after header"),
		},
		{
			name:          "Unexpected Status Code",
			statusCode:    http.StatusInternalServerError,
			expectedError: errors.New("unexpected status code: 500 Internal Server Error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.statusCode)
				if tt.retryAfter != "" {
					w.Header().Set("Retry-After", tt.retryAfter)
				}
			}))
			defer server.Close()

			_, err := GetAccrualOrderStatus(server.URL, "123")

			assert.Error(t, err)
		})
	}
}
