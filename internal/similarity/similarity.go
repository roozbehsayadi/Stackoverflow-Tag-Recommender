package similarity

import (
	"sync"
)

var questions map[int]map[string]bool

func InitData(questionsIn map[int]map[string]bool) {
	questions = questionsIn
}

func GetSimilarities(tags chan []string, result chan TagSimilarity, wg *sync.WaitGroup) {
	defer wg.Done()
	for currentTags := range tags {
		tag1 := currentTags[0]
		tag2 := currentTags[1]
		var either, both int
		for _, tags := range questions {
			if tags[tag1] && tags[tag2] {
				both++
			}
			if tags[tag1] || tags[tag2] { // xor
				either++
			}
		}
		result <- TagSimilarity{
			Tag1:       tag1,
			Tag2:       tag2,
			Similarity: (float64(both) / float64(either)),
		}
	}
}

type TagSimilarity struct {
	Tag1       string
	Tag2       string
	Similarity float64
}
