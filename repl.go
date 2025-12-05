package main

import (
	"strings"
)

func cleanInput(text string) []string {
	rawWords := strings.Split(text, " ")
	var words []string

	for i := range len(rawWords) {
		if rawWords[i] == "" {
			continue
		}
		word := strings.Trim(rawWords[i], " ")
		word = strings.ToLower(word)
		words = append(words, word)
	}

	return words
}
