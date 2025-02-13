package domain

type Balance struct {
	Current  int `json:"current"`
	Withdraw int `json:"withdraw"`
}
