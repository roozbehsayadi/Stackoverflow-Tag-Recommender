package main

import (
	"context"
	"fmt"
	"math/rand"
	"sort"
	"stackoverflow-recommender/internal"
	"stackoverflow-recommender/internal/postgres"
	"stackoverflow-recommender/internal/similarity"
	"sync"
)

const pgxUrl = "postgres://postgres:postgres@localhost:5432/stackoverflow_recommender"

func main() {

	conn := postgres.GetDatabaseConn(pgxUrl)
	defer conn.Close(context.Background())

	questions, allTags := postgres.RetrieveDataFromDatabase(conn, 100)

	similarity.InitData(questions)

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

	results := []similarity.TagSimilarity{}
	appendResultsWaitGroup.Add(1)
	go func() {
		defer appendResultsWaitGroup.Done()
		for result := range res {
			results = append(results, result)
			if rand.Int()%1000 == 0 {
				fmt.Printf("%f%%\n", float64(len(results))/float64(
					len(allTagsAr)*(len(allTagsAr)-1)/2,
				))
			}
		}
	}()
	calculationsWaitGroup.Wait()
	close(res)
	appendResultsWaitGroup.Wait()
	sort.Slice(results, func(i, j int) bool {
		return results[i].Similarity > results[j].Similarity
	})

	err := internal.WriteResultsToFile("similarities.csv", results)
	if err != nil {
		fmt.Println("Could not write to file. Writing here...")
		internal.WriteResultsToStdout(results)
	}

}
