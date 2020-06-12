package v4

import (
	"context"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
)

type Querier interface {
	QueryRow(ctx context.Context, query string, args ...interface{}) pgx.Row
	Query(ctx context.Context, query string, args ...interface{}) (pgx.Rows, error)
}

type Executer interface {
	Exec(ctx context.Context, query string, args ...interface{}) (pgconn.CommandTag, error)
}

type ExecQueryer interface {
	Executer
	Querier
}

type Beginner interface {
	Begin(ctx context.Context) (pgx.Tx, error)
}

type Transactioner interface {
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
}

type DriverWrapper interface {
	Beginner
	ExecQueryer
}
