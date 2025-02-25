package jobs

import (
	"encoding/json"
	"errors"
	"fmt"
	"gophermart/internal/domain"
	"net/http"
	"strconv"

	"github.com/go-resty/resty/v2"
)

func (j *OrdersJob) GetOrderStatus(order string) (*domain.AccrualResponse, error) {
	client := resty.New()
	uri := j.Utils.C.AccrualHost + "/api/orders/" + order
	resp, err := client.R().Get(uri)
	if err != nil {
		return nil, fmt.Errorf("failed to get order status: %w", err)
	}
	switch resp.StatusCode() {
	case http.StatusOK:
		var orderResponse domain.AccrualResponse
		if err := json.Unmarshal(resp.Body(), &orderResponse); err != nil {
			return nil, err
		}
		return &orderResponse, nil
	case http.StatusTooManyRequests:
		retryAfter := resp.Header().Get("Retry-After")
		if retryAfter == "" {
			return nil, errors.New("no retry-after header")
		}
		if seconds, err := strconv.Atoi(retryAfter); err == nil {
			return nil, &domain.TooManyRequestsError{RetryAfter: seconds}
		}
		return nil, errors.New("invalid retry-after header")
	default:
		return nil, errors.New("unexpected status code" + resp.Status())
	}
}
