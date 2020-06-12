package db_forum

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
)

func ConnectToDB(dbHost, dbName, dbUser, dbPassword string, dbMaxConns int) (*pgxpool.Pool, error) {
	connStr := fmt.Sprintf(
		"host=%s dbname=%s user=%s password=%s pool_max_conns=%d",
		dbHost, dbName, dbUser, dbPassword, dbMaxConns,
	)

	return pgxpool.Connect(context.Background(), connStr)

	//return pgx.NewConnPool(
	//	pgx.ConnPoolConfig{
	//		ConnConfig: pgx.ConnConfig{
	//			Host:           dbHost,
	//			Port:           5432,
	//			Database:       dbName,
	//			User:           dbUser,
	//			Password:       dbPassword,
	//			TLSConfig:      nil,
	//			UseFallbackTLS: false,
	//		},
	//		MaxConnections: 10,
	//		AfterConnect:   nil,
	//		AcquireTimeout: 0,
	//	},
	//)
}
