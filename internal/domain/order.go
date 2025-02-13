package domain

import "time"

type OrderWithUserID struct {
	Number string `json:"number"`
	UserID int    `json:"-"`
}

type Order struct {
	OrderWithUserID
	Accrual    float64   `json:"accrual,omitempty"`
	Status     string    `json:"status"`
	UploadedAt time.Time `json:"uploaded_at"`
}
