package users

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"gophermart/internal/domain"
)

type UserRepo struct {
	conn *pgxpool.Pool
}

func New(conn *pgxpool.Pool) *UserRepo {
	return &UserRepo{conn: conn}
}

const InsertUser = "INSERT INTO users (login,password) VALUES ($1,$2)"

func (u *UserRepo) SaveUser(ctx context.Context, values *domain.Credentials) (err error) {
	_, err = u.conn.Exec(ctx, InsertUser, values.Login, values.Password)
	return err
}

const SelectUser = "SELECT id, password FROM users WHERE login = $1"

func (u *UserRepo) GetUser(ctx context.Context, login string) (*domain.UserIDPassword, error) {
	var user domain.UserIDPassword
	err := u.conn.QueryRow(ctx, SelectUser, login).Scan(&user.ID, &user.Password)
	return &user, err
}
