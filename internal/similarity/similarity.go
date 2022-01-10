package similarity

import (
	"sync"
)

var tagsToQuestions map[string]map[int]bool

func InitData(tags map[string]map[int]bool) {
	tagsToQuestions = tags
}

func GetSimilarities(tags chan []string, result chan TagSimilarity, processedAnAnswer chan bool, wg *sync.WaitGroup) {
	defer wg.Done()
	for currentTags := range tags {
		tag1 := currentTags[0]
		tag2 := currentTags[1]
		if (len(tagsToQuestions[tag1]) < 5 || len(tagsToQuestions[tag2]) < 5) &&
			!MustStoreResult(tag1, tag2) {
			processedAnAnswer <- true
			continue
		}
		var either, both int
		for question := range tagsToQuestions[tag1] {
			if tagsToQuestions[tag2][question] {
				both++
			}
			either++
		}
		for question := range tagsToQuestions[tag2] {
			if !tagsToQuestions[tag1][question] {
				either++
			}
		}
		similarity := (float64(both) / float64(either))
		if similarity < 0.3 && !MustStoreResult(tag1, tag2) {
			processedAnAnswer <- true
			continue
		}
		processedAnAnswer <- true
		result <- TagSimilarity{
			Tag1:       tag1,
			Tag2:       tag2,
			Similarity: similarity,
		}
	}
}

type TagSimilarity struct {
	Tag1       string
	Tag2       string
	Similarity float64
}
