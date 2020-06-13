package sql

type Scanner interface {
	Scan(dst ...interface{}) error
}
