package main

import (
	"fmt"
	"testing"
)

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input    string
		expected []string
	}{
		{
			input:    "   hello world  ",
			expected: []string{"hello", "world"},
		},
		{
			input:    "ThIs Is Soome Werid   casing  ",
			expected: []string{"this", "is", "soome", "werid", "casing"},
		},
	}

	for _, c := range cases {
		actual := cleanInput(c.input)

		fmt.Printf("Expecting: %v\nRecieved: %v\n", c.expected, actual)

		if len(actual) != len(c.expected) {
			t.Errorf("Length of function output doesnt match expected")
		}

		for i := range actual {
			word := actual[i]
			expectedWord := c.expected[i]

			if word != expectedWord {
				t.Errorf("Words do not match at index: %v | %v != %v", i, word, expectedWord)
			}
		}
	}
}
