package internal

import (
	"fmt"
	"os"
	"stackoverflow-recommender/internal/similarity"
)

func WriteResultsToFile(address string, results []similarity.TagSimilarity) error {
	f, err := os.OpenFile("similarities.csv", os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.Write([]byte("tag1,tag2,similarity\n"))
	if err != nil {
		return err
	}
	for _, result := range results {
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
