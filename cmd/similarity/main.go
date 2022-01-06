package main

import (
	"context"
	"fmt"
	"os"
	"stackoverflow-recommender/internal/postgres"
)

const pgxUrl = "postgres://postgres:postgres@localhost:5432/stackoverflow_recommender"

func main() {

	conn := postgres.GetDatabaseConn(pgxUrl)
	defer conn.Close(context.Background())

	rows, err := conn.Query(context.Background(), "SELECT question_id, tag FROM question_tags LIMIT 10")
	if err != nil {
		fmt.Println("Query was not sucessfull: ", err)
		os.Exit(1)
	}
	var question_id int
	var tag string

	for rows.Next() {
		rows.Scan(&question_id, &tag)
		fmt.Println(question_id, tag)
	}

}
