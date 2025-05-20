package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
)

func main() {
	log.SetFlags(0)

	filenames := os.Args[1:]
	totalWordCount := 0
	didError := false

	for _, filename := range filenames {
		wordcount, err := CountWordsInFile(filename)
		if err != nil {
			didError = true
			fmt.Printf("%s: %s\n", filename, err)
			continue
		}
		totalWordCount += wordcount
		fmt.Printf("%s: %d\n", filename, wordcount)
	}

	if len(filenames) == 0 {
		wordCount := CountWords(os.Stdin)
		fmt.Println(wordCount)
	}

	if len(filenames) > 0 {
		fmt.Printf("Total word count: %d\n", totalWordCount)
	}

	if didError {
		os.Exit(1)
	}
}

func CountWordsInFile(filename string) (int, error) {
	file, err := os.Open(filename)
	if err != nil {
		return 0, fmt.Errorf("open file: %w", err)
	}
	defer file.Close()

	wordCount := CountWords(file)

	return wordCount, nil
}

func CountWords(file io.Reader) int {
	wordCount := 0

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanWords)

	for scanner.Scan() {
		wordCount++
	}

	return wordCount
}
