package similarity

import (
	"sync"
)

var questionsToTags map[int]map[string]bool
var tagsToQuestions map[string]map[int]bool

func InitData(questions map[int]map[string]bool, tags map[string]map[int]bool) {
	questionsToTags = questions
	tagsToQuestions = tags
}

func GetSimilarities(tags chan []string, result chan TagSimilarity, wg *sync.WaitGroup) {
	defer wg.Done()
	for currentTags := range tags {
		tag1 := currentTags[0]
		tag2 := currentTags[1]
		var either, both int
		for question, _ := range tagsToQuestions[tag1] {
			if tagsToQuestions[tag2][question] {
				both++
			}
			either++
		}
		for question, _ := range tagsToQuestions[tag2] {
			if !tagsToQuestions[tag1][question] {
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
