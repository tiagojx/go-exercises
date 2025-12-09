package main

import (
	"fmt"
	"strings"
)

func main() {
	fmt.Println("EX01: Word Counter")

	// Split the sentence into a array of words.
	sentence := "the quick brown fox jumps over the lazy dog dog"
	sentenceSplited := strings.Split(sentence, " ")

	// Make a empty Go's map.
	wordCounts := make(map[string]int)

	for word := range len(sentenceSplited) {
		wordCounts[sentenceSplited[word]] += 1
	}

	for word, count := range wordCounts {
		fmt.Printf("The word '%s' appears %d times.\n", word, count)
	}
}
