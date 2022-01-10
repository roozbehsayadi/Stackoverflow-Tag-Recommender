package internal

import (
	"fmt"
	"sort"
	"stackoverflow-recommender/internal/file"
	"stackoverflow-recommender/internal/similarity"
)

var requiredTags map[string]bool

func init() {
	requiredTags = make(map[string]bool)
	requiredTags["intellij-idea"] = true
	requiredTags["jax-rs"] = true
	requiredTags["user-interface"] = true
	requiredTags["regex"] = true
	requiredTags["static"] = true
	requiredTags["session"] = true
	requiredTags["spring"] = true
	requiredTags["nullpointerexception"] = true
	requiredTags["dependency-injection"] = true
}

func GenerateReports(results []similarity.TagSimilarity) {
	fmt.Println("Generating report for Q1")
	generateQ1Report(results)
	fmt.Println("Generating report for Q2")
	generateQ2Report(results)
	fmt.Println("Generating report for Q3")
	generateQ3Report(results)
}

func generateQ1Report(results []similarity.TagSimilarity) {

	intellijSimilarTags := getTop10SimilarTagSimilarities(results, "intellij-idea")
	jaxSimilarTags := getTop10SimilarTagSimilarities(results, "jax-rs")
	userInterfaceSimilarTags := getTop10SimilarTagSimilarities(results, "user-interface")

	similarTags := []string{}
	similarTags = append(similarTags, getSimilarTags("intellij-idea", intellijSimilarTags)...)
	similarTags = append(similarTags, getSimilarTags("jax-rs", jaxSimilarTags)...)
	similarTags = append(similarTags, getSimilarTags("user-interface", userInterfaceSimilarTags)...)

	file.WriteArrayTofile("q1_report.csv", similarTags)

}

func generateQ2Report(results []similarity.TagSimilarity) {
	tagPairs := [][]string{{"regex", "static"}, {"session", "spring"}, {"nullpointerexception", "dependency-injection"}}
	outputArray := []string{"Tag1,Tag2,distance", "\n"}
	for _, tagSimilarity := range results {
		if matchPairs(tagSimilarity, tagPairs) {
			outputArray = append(outputArray, tagSimilarity.Tag1+","+
				tagSimilarity.Tag2+","+
				fmt.Sprintf("%f", 1-tagSimilarity.Similarity), "\n",
			)
		}
	}

	file.WriteArrayTofile("q2_report.csv", outputArray)
}

func generateQ3Report(results []similarity.TagSimilarity) {
	outputArray := []string{"Tag1,Tag2,similarity", "\n"}
	for _, tagSimilarity := range results[0:10] {
		outputArray = append(outputArray, tagSimilarity.Tag1+","+
			tagSimilarity.Tag2+","+
			fmt.Sprintf("%f", tagSimilarity.Similarity), "\n",
		)
	}

	file.WriteArrayTofile("q3_report.csv", outputArray)
}

func matchPairs(tagSimilarity similarity.TagSimilarity, tagPairs [][]string) bool {
	for _, pair := range tagPairs {
		if (tagSimilarity.Tag1 == pair[0] && tagSimilarity.Tag2 == pair[1]) ||
			(tagSimilarity.Tag2 == pair[0] && tagSimilarity.Tag1 == pair[1]) {
			return true
		}
	}
	return false
}

func getSimilarTags(tag string, similarities []similarity.TagSimilarity) []string {
	mostSimilars := []string{}
	mostSimilars = append(mostSimilars, tag+":")
	for _, tagSimilarity := range similarities {
		var similarTag string
		if tagSimilarity.Tag1 == tag {
			similarTag = tagSimilarity.Tag2
		} else {
			similarTag = tagSimilarity.Tag1
		}
		mostSimilars = append(mostSimilars, similarTag)
	}
	mostSimilars = append(mostSimilars, "\n")
	return mostSimilars
}

func getTop10SimilarTagSimilarities(similarities []similarity.TagSimilarity, tag string) []similarity.TagSimilarity {
	allSimilarities := []similarity.TagSimilarity{}
	for _, tagSimilarity := range similarities {
		if tagSimilarity.Tag1 == tag || tagSimilarity.Tag2 == tag {
			allSimilarities = append(allSimilarities, tagSimilarity)
		}
	}
	sort.Slice(allSimilarities, func(i, j int) bool {
		return allSimilarities[i].Similarity > allSimilarities[j].Similarity
	})
	if len(allSimilarities) < 10 {
		return allSimilarities
	} else {
		return allSimilarities[0:10]
	}
}
