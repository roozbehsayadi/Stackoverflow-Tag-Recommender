package postgres

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v4"
)

func GetDatabaseConn(url string) *pgx.Conn {
	conn, err := pgx.Connect(context.Background(), url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	return conn
}
