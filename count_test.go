package main_test

import (
	"strings"
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
			reader := strings.NewReader(tc.input)
			result := counter.CountWords(reader)
			if result != tc.expected {
				t.Logf("Expected %d, but got %d", tc.expected, result)
				t.Fail()
			}
		})
	}
}

func TestCountLines(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected int
	}{
		{"1 line no newline", "one two three", 0},
		{"Empty file", "", 0},
		{"2 lines with 1 newline", "one two three\nfour five", 1},
		{"2 line with 2 newlines", "one two three\nfour five\n", 2},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			reader := strings.NewReader(tc.input)
			result := counter.CountLines(reader)
			if result != tc.expected {
				t.Logf("Expected %d, but got %d", tc.expected, result)
				t.Fail()
			}
		})
	}
}

func TestCountBytes(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected int
	}{
		{"Empty file", "", 0},
		{"1 byte", "a", 1},
		{"2 bytes", "ab", 2},
		{"3 bytes", "abc", 3},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			reader := strings.NewReader(tc.input)
			result := counter.CountBytes(reader)
			if result != tc.expected {
				t.Logf("Expected %d, but got %d", tc.expected, result)
				t.Fail()
			}
		})
	}
}

func TestGetCounts(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected counter.Counts
	}{
		{
			"simple five words",
			"one two three four five\n",
			counter.Counts{Lines: 1, Words: 5, Bytes: 24},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			reader := strings.NewReader(tc.input)
			result := counter.GetCounts(reader)
			if result != tc.expected {
				t.Logf("Expected %v, but got %v", tc.expected, result)
				t.Fail()
			}
		})
	}
}

var benchData = []string{
	"This is a test string\nthat spans\nmultiple lines",
	"one two three four\n five six\n seven eight nine\n \n ten",
	"Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed euismod, nisl eget aliquam aliquet, nunc nisl aliquet nisl, eget aliquet nisl nisl eget nisl. Sed euismod, nisl eget aliquam aliquet, nunc nisl aliquet nisl, eget aliquet nisl nisl eget nisl.",
}

func BenchmarkGetCounts(b *testing.B) {
	for i := range b.N {
		data := benchData[i%len(benchData)]
		reader := strings.NewReader(data)
		counter.GetCounts(reader)
	}
}

func BenchmarkGetCountsSinglePass(b *testing.B) {
	for i := range b.N {
		data := benchData[i%len(benchData)]
		reader := strings.NewReader(data)
		counter.GetCountsSinglePass(reader)
	}
}
