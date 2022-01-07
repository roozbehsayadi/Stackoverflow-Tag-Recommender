package main

import (
	"context"
	"fmt"
	"stackoverflow-recommender/internal"
	"stackoverflow-recommender/internal/postgres"
	"stackoverflow-recommender/internal/similarity"
	"sync"
)

const pgxUrl = "postgres://postgres:postgres@localhost:5432/stackoverflow_recommender"

func main() {

	conn := postgres.GetDatabaseConn(pgxUrl)
	defer conn.Close(context.Background())
	fmt.Println("Connected to database")

	questionsToTags, tagsToQuestions, allTags := postgres.RetrieveDataFromDatabase(conn, -1)
	fmt.Println("Retrieved and processed data from database")

	similarity.InitData(questionsToTags, tagsToQuestions)

	var allTagsAr []string
	for key := range allTags {
		allTagsAr = append(allTagsAr, key)
	}

	allTagsChan := make(chan ([]string))

	go func() {
		defer close(allTagsChan)
		for i := range allTagsAr {
			for j := i + 1; j < len(allTagsAr); j++ {
				allTagsChan <- []string{allTagsAr[i], allTagsAr[j]}
			}
		}
	}()

	res := make(chan (similarity.TagSimilarity))
	var calculationsWaitGroup, appendResultsWaitGroup sync.WaitGroup

	for i := 0; i < 8; i++ {
		calculationsWaitGroup.Add(1)
		go similarity.GetSimilarities(allTagsChan, res, &calculationsWaitGroup)
	}

	// results := []similarity.TagSimilarity{}
	appendResultsWaitGroup.Add(1)
	go func() {
		defer appendResultsWaitGroup.Done()
		internal.WriteResultsToFile("similarities.csv", res, len(allTagsAr)*(len(allTagsAr)-1)/2)
	}()
	calculationsWaitGroup.Wait()
	close(res)
	appendResultsWaitGroup.Wait()

}
