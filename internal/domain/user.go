package domain

type Credentials struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type UserIDPassword struct {
	ID       int    `json:"id"`
	Password string `json:"password"`
}
