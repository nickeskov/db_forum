package sql

import "github.com/jackc/pgx"

type Querier interface {
	QueryRow(query string, args ...interface{}) *pgx.Row
	Query(query string, args ...interface{}) (*pgx.Rows, error)
}

type Executer interface {
	Exec(query string, args ...interface{}) (pgx.CommandTag, error)
}

type ExecQueryer interface {
	Executer
	Querier
}

type Beginner interface {
	Begin() (*pgx.Tx, error)
}

type DriverWrapper interface {
	Beginner
	ExecQueryer
}
