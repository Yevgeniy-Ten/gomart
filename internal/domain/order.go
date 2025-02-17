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
type OrderWithAccrual struct {
	Number  string      `json:"number"`
	Accrual *float64    `json:"accrual"`
	Status  OrderStatus `json:"status"`
}
