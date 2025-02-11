package repository

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose"
	"gophermart/internal/domain"
)

type Database struct {
	conn *pgxpool.Pool
}

func (d *Database) Close(_ context.Context) {
	d.conn.Close()
}

func New(utils *domain.Utils) (*Database, error) {
	c := context.TODO()
	conn, err := pgxpool.New(c, utils.C.DatabaseURL)
	if err != nil {
		return nil, err
	}
	d := &Database{
		conn: conn,
	}
	if err := d.Init(); err != nil {
		return nil, err
	}

	return d, nil
}

func (d *Database) Init() error {
	db := stdlib.OpenDBFromPool(d.conn)
	if err := goose.Up(db, "./migrations"); err != nil {
		return err
	}
	return nil
}
