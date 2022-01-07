package internal

import (
	"fmt"
	"math/rand"
	"os"
	"stackoverflow-recommender/internal/similarity"
)

func WriteResultsToFile(address string, results chan similarity.TagSimilarity, expectedNumberOfResults int) error {
	f, err := os.OpenFile("similarities.csv", os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.Write([]byte("tag1,tag2,similarity\n"))
	if err != nil {
		return err
	}
	counter := 0
	for result := range results {
		counter++
		if rand.Int()%5000 == 0 {
			fmt.Printf("%f%%\n", 100*float64(counter)/float64(expectedNumberOfResults))
		}
		if result.Similarity < 0.3 {
			continue
		}
		temp := result.Tag1 + "," + result.Tag2 + "," + fmt.Sprintf("%f", result.Similarity) + "\n"
		_, err := f.Write([]byte(temp))
		if err != nil {
			return err
		}
	}
	return nil
}

func WriteResultsToStdout(results []similarity.TagSimilarity) {
	fmt.Println("tag1,tag2,similarity")
	for _, result := range results {
		fmt.Printf("%v,%v,%f\n", result.Tag1, result.Tag2, result.Similarity)
	}
}
