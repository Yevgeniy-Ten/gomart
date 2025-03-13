package repository

import (
	"context"
	"gophermart/internal/domain"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose"
)

//go:generate mockgen -source=repository.go -destination=mocks/mock_dbpool.go -package=mocks
type DBPool interface {
	Exec(ctx context.Context, sql string, args ...interface{}) (pgconn.CommandTag, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	Close()
}

type Repo struct {
	conn DBPool
	pool *pgxpool.Pool
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
	d := NewWithPool(conn)
	d.pool = conn
	return d, nil
}
func NewWithPool(pool DBPool) *Repo {
	return &Repo{conn: pool}
}
func (d *Repo) Init() error {
	db := stdlib.OpenDBFromPool(d.pool)
	if err := goose.Up(db, "./migrations"); err != nil {
		return err
	}
	return nil
}
