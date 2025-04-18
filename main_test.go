package main_test

import (
	"testing"

	counter "github.com/ColeBurch/counter"
)

func TestCountWords(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected int
	}{
		{"5 words", "one two three four five", 5},
		{"Empty file", "", 0},
		{"Single space", " ", 0},
		{"new lines", "one two three\nfour five", 5},
		{"multiple spaces", "one  two   three    four five", 5},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := counter.CountWords([]byte(tc.input))
			if result != tc.expected {
				t.Logf("Expected %d, but got %d", tc.expected, result)
				t.Fail()
			}
		})
	}
}
