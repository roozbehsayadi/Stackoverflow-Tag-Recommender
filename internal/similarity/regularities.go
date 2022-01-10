package similarity

type tagPair struct {
	tag1 string
	tag2 string
}

var requiredTags map[string]bool
var requiredPairs map[tagPair]bool

func init() {
	requiredTags = make(map[string]bool)
	requiredPairs = make(map[tagPair]bool)
	requiredTags["intellij-idea"] = true
	requiredTags["jax-rs"] = true
	requiredTags["user-interface"] = true

	requiredPairs[tagPair{tag1: "regex", tag2: "static"}] = true
	requiredPairs[tagPair{tag1: "session", tag2: "spring"}] = true
	requiredPairs[tagPair{tag1: "nullpointerexception", tag2: "dependency-injection"}] = true
}

func MustStoreResult(tag1 string, tag2 string) bool {
	if requiredTags[tag1] || requiredTags[tag2] {
		return true
	}
	temp1 := tagPair{tag1: tag1, tag2: tag2}
	if requiredPairs[temp1] {
		return true
	}
	temp2 := tagPair{tag1: tag2, tag2: tag1}
	return requiredPairs[temp2]
}
