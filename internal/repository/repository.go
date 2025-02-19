package repository

import (
	"context"
	"gophermart/internal/domain"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose"
)

type Repo struct {
	conn *pgxpool.Pool
}

func (d *Repo) Close(_ context.Context) {
	d.conn.Close()
}

func New(utils *domain.Utils) (*Repo, error) {
	c := context.TODO()
	conn, err := pgxpool.New(c, utils.C.DatabaseURL)
	if err != nil {
		return nil, err
	}
	d := &Repo{
		conn: conn,
	}
	if err := d.Init(); err != nil {
		return nil, err
	}

	return d, nil
}

func (d *Repo) Init() error {
	db := stdlib.OpenDBFromPool(d.conn)
	if err := goose.Up(db, "./migrations"); err != nil {
		return err
	}
	return nil
}
