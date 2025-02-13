package repository

import (
	"context"
)

const CreateBalance = "INSERT INTO balance (user_id) VALUES ($1)"

func (d *Repo) CreateBalance(userID int) error {
	_, err := d.conn.Exec(context.Background(), CreateBalance, userID)
	return err
}
