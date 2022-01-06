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

func RetrieveDataFromDatabase(conn *pgx.Conn, rowsCount int) (map[int]map[string]bool, map[string]bool) {
	rows, err := conn.Query(context.Background(), "SELECT question_id, tag FROM question_tags ORDER BY question_id LIMIT $1", rowsCount)
	if err != nil {
		fmt.Println("Query was not sucessfull:", err)
		os.Exit(1)
	}

	questions := make(map[int]map[string]bool)
	allTags := make(map[string]bool)

	var question_id int
	var tag string

	for rows.Next() {
		rows.Scan(&question_id, &tag)
		allTags[tag] = true
		if _, ok := questions[question_id]; !ok {
			questions[question_id] = make(map[string]bool)
		}
		questions[question_id][tag] = true
	}
	return questions, allTags
}
