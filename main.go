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

	file, err := os.Open("words.txt")
	if err != nil {
		log.Fatalln("Error reading file:", err)
	}

	wordCount := CountWords(file)
	fmt.Println("Word Count: " + fmt.Sprint(wordCount))
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
