package sql

type Transactioner interface {
	Commit() error
	Rollback() error
}
