package domain

import "time"

type Balance struct {
	Current  float64 `json:"current"`
	Withdraw float64 `json:"withdraw"`
}
type OrderToWithdraw struct {
	Order string  `json:"order"`
	Sum   float64 `json:"sum"`
}
type Withdraw struct {
	OrderToWithdraw
	ProcessedAt time.Time `json:"processed_at"`
}
