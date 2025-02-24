package domain

import "time"

type OrderStatus string

const (
	OrderStatusNew        OrderStatus = "NEW"
	OrderStatusProcessed  OrderStatus = "PROCESSED"
	OderStatusInvalid     OrderStatus = "INVALID"
	OrderStatusProcessing OrderStatus = "PROCESSING"
)

type OrderWithUserID struct {
	Number string `json:"number"`
	UserID int    `json:"-"`
}

type Order struct {
	OrderWithUserID
	Accrual    float64     `json:"accrual,omitempty"`
	Status     OrderStatus `json:"status"`
	UploadedAt time.Time   `json:"uploaded_at"`
}
type AccrualResponse struct {
	Status  OrderStatus `json:"status"`
	Accrual *float64    `json:"accrual,omitempty"`
}

type OrderWithAccrual struct {
	AccrualResponse
	OrderWithUserID
}
type OrderInJobs struct {
	Number string `json:"number"`
	AccrualResponse
	Error  error
	UserID int
}

type TooManyRequestsError struct {
	RetryAfter int
}

func (e *TooManyRequestsError) Error() string {
	return "too many requests"
}
