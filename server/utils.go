package server

import (
	"regexp"
	"strings"
)

var (
	wordRE                    = regexp.MustCompile(`^[a-z]+$`)
	metadataKeyRE             = regexp.MustCompile(`^[a-zA-Z]+$`)
	defaultMetadataValueKeyRE = regexp.MustCompile(`^[a-z]+/[a-zA-Z]+$`)
)

func splitJoinedWords(s string) []string {
	words := strings.Split(s, ",")
	if len(words) == 0 {
		return nil
	}
	for _, w := range words {
		if !isValidWord(w) {
			return nil
		}
	}
	return words
}

func isValidTag(tag string) bool {
	return isValidWord(tag)
}

func isValidCollectionKey(key string) bool {
	return isValidWord(key)
}

func isValidWord(word string) bool {
	return wordRE.MatchString(word)
}

func isValidMetadataKey(key string) bool {
	return metadataKeyRE.MatchString(key)
}

func isValidDefaultMetadataValueKeyRE(key string) bool {
	return defaultMetadataValueKeyRE.MatchString(key)
}
