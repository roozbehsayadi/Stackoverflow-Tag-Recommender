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

var (
	dbUsername = "postgres"
	dbPassword = "postgres"
	dbHost     = "localhost"
	dbPort     = "5432"
	dbName     = "stackoverflow_recommender"
	pgxUrl     = "postgres://" + dbUsername + ":" + dbPassword + "@" + dbHost + ":" + dbPort + "/" + dbName
)

func main() {

	conn := postgres.GetDatabaseConn(pgxUrl)
	defer conn.Close(context.Background())
	fmt.Println("Connected to database")

	tagsToQuestions, allTags := postgres.RetrieveDataFromDatabase(conn, -1)
	fmt.Println("Retrieved and processed data from database")

	similarity.InitData(tagsToQuestions)

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
	processedAnAnswer := make(chan bool)
	var calculationsWaitGroup, appendResultsWaitGroup sync.WaitGroup

	for i := 0; i < 8; i++ {
		calculationsWaitGroup.Add(1)
		go similarity.GetSimilarities(allTagsChan, res, processedAnAnswer, &calculationsWaitGroup)
	}

	results := []similarity.TagSimilarity{}
	appendResultsWaitGroup.Add(1)
	go func() {
		defer appendResultsWaitGroup.Done()
		for result := range res {
			results = append(results, result)
		}
		close(processedAnAnswer)
	}()
	appendResultsWaitGroup.Add(1)
	go func() {
		defer appendResultsWaitGroup.Done()
		counter := 0
		expectedNumberOfResults := len(allTagsAr) * (len(allTagsAr) - 1) / 2
		for range processedAnAnswer {
			counter++
			if rand.Int()%5000 == 0 {
				fmt.Printf("%f%%\n", 100*float64(counter)/float64(expectedNumberOfResults))
			}
		}
		fmt.Println(100*float64(counter)/float64(expectedNumberOfResults), "%")
	}()
	calculationsWaitGroup.Wait()
	close(res)
	appendResultsWaitGroup.Wait()

	sort.Slice(results, func(i, j int) bool {
		return results[i].Similarity > results[j].Similarity
	})

	internal.GenerateReports(results)

}
