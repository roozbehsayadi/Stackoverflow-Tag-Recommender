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

func RetrieveDataFromDatabase(conn *pgx.Conn, rowsCount int) (map[int]map[string]bool, map[string]map[int]bool, map[string]bool) {
	var rows pgx.Rows
	var err error
	if rowsCount == -1 {
		rows, err = conn.Query(context.Background(), "SELECT question_id, tag FROM question_tags")
	} else {
		rows, err = conn.Query(context.Background(), "SELECT question_id, tag FROM question_tags ORDER BY question_id LIMIT $1", rowsCount)
	}
	if err != nil {
		fmt.Println("Query was not sucessfull:", err)
		os.Exit(1)
	}

	questionsToTags := make(map[int]map[string]bool)
	tagsToQuestions := make(map[string]map[int]bool)
	allTags := make(map[string]bool)

	var question_id int
	var tag string

	for rows.Next() {
		rows.Scan(&question_id, &tag)
		allTags[tag] = true
		if _, ok := questionsToTags[question_id]; !ok {
			questionsToTags[question_id] = make(map[string]bool)
		}
		questionsToTags[question_id][tag] = true
		if _, ok := tagsToQuestions[tag]; !ok {
			tagsToQuestions[tag] = make(map[int]bool)
		}
		tagsToQuestions[tag][question_id] = true
	}
	return questionsToTags, tagsToQuestions, allTags
}
