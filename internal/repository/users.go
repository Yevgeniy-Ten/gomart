package repository

import (
	"context"
	"fmt"
	"gophermart/internal/domain"
)

const InsertUser = "INSERT INTO users (login,password) VALUES ($1,$2) RETURNING id"

func (d *Repo) SaveUser(ctx context.Context, values *domain.Credentials) error {
	tx, err := d.conn.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	var userID int
	err = tx.QueryRow(ctx, InsertUser, values.Login, values.Password).Scan(&userID)
	if err != nil {
		return fmt.Errorf("failed to insert user: %w", err)
	}

	_, err = tx.Exec(ctx, CreateBalance, userID)
	if err != nil {
		return fmt.Errorf("failed to create balance: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

const SelectUser = "SELECT id, password FROM users WHERE login = $1"

func (d *Repo) GetUser(ctx context.Context, login string) (*domain.UserIDPassword, error) {
	var user domain.UserIDPassword
	err := d.conn.QueryRow(ctx, SelectUser, login).Scan(&user.ID, &user.Password)
	return &user, err
}
