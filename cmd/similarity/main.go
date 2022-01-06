package main

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"sort"
	"stackoverflow-recommender/internal/postgres"
	"stackoverflow-recommender/internal/similarity"
	"sync"
)

const pgxUrl = "postgres://postgres:postgres@localhost:5432/stackoverflow_recommender"

func main() {

	conn := postgres.GetDatabaseConn(pgxUrl)
	defer conn.Close(context.Background())

	rows, err := conn.Query(context.Background(), "SELECT question_id, tag FROM question_tags ORDER BY question_id LIMIT 500")
	if err != nil {
		fmt.Println("Query was not sucessfull: ", err)
		os.Exit(1)
	}
	var question_id int
	var tag string

	questions := make(map[int]map[string]bool)
	allTags := make(map[string]bool)

	for rows.Next() {
		rows.Scan(&question_id, &tag)
		allTags[tag] = true
		if _, ok := questions[question_id]; !ok {
			questions[question_id] = make(map[string]bool)
		}
		questions[question_id][tag] = true
	}

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
	var wg, wg2 sync.WaitGroup

	for i := 0; i < 8; i++ {
		wg.Add(1)
		go similarity.GetSimilarities(allTagsChan, res, &wg)
	}
	results := []similarity.TagSimilarity{}
	wg2.Add(1)
	go func() {
		defer wg2.Done()
		for result := range res {
			results = append(results, result)
			if rand.Int()%1000 == 0 {
				fmt.Printf("%f%%\n", float64(len(results))/float64(
					len(allTagsAr)*(len(allTagsAr)-1)/2,
				))
			}
		}
	}()
	wg.Wait()
	close(res)
	wg2.Wait()
	sort.Slice(results, func(i, j int) bool {
		return results[i].Similarity < results[j].Similarity
	})

	f, err := os.OpenFile("similarities.csv", os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Could not create the file.")
		os.Exit(1)
	}
	defer f.Close()
	_, err = f.Write([]byte("tag1,tag2,similarity\n"))
	if err != nil {
		fmt.Println("Error in writing to file:", err)
		os.Exit(1)
	}
	for _, result := range results {
		temp := result.Tag1 + "," + result.Tag2 + "," + fmt.Sprintf("%f", result.Similarity) + "\n"
		_, err := f.Write([]byte(temp))
		if err != nil {
			fmt.Println("Error in writing to file:", err)
			os.Exit(1)
		}
	}
	// for _, result := range results {
	// 	fmt.Println(result.Tag1, result.Tag2, result.Similarity)
	// }
}
