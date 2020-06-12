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
}
