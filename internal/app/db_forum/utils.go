package db_forum

import "github.com/jackc/pgx"

func ConnectToDB(dbHost, dbName, dbUser, dbPassword string) (*pgx.ConnPool, error) {
	return pgx.NewConnPool(
		pgx.ConnPoolConfig{
			ConnConfig: pgx.ConnConfig{
				Host:           dbHost,
				Port:           5432,
				Database:       dbName,
				User:           dbUser,
				Password:       dbPassword,
				TLSConfig:      nil,
				UseFallbackTLS: false,
			},
			MaxConnections: 10,
			AfterConnect:   nil,
			AcquireTimeout: 0,
		},
	)
}
