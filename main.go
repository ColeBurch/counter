package main

import (
	"bytes"
	"fmt"
	"os"
)

func main() {
	data, _ := os.ReadFile("words.txt")
	wordCount := CountWords(data)
	fmt.Println("Word Count: " + fmt.Sprint(wordCount))
}

func CountWords(data []byte) int {
	words := bytes.Fields(data)
	return len(words)
}
