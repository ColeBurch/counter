package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

type Counts struct {
	Lines int
	Words int
	Bytes int
}

func GetCounts(file io.ReadSeeker) Counts {
	const offsetStart = int64(0)

	lineCount := CountLines(file)
	file.Seek(offsetStart, io.SeekStart)
	wordCount := CountWords(file)
	file.Seek(offsetStart, io.SeekStart)
	byteCount := CountBytes(file)

	return Counts{
		Lines: lineCount,
		Words: wordCount,
		Bytes: byteCount,
	}
}

func CountFile(filename string) (Counts, error) {
	file, err := os.Open(filename)
	if err != nil {
		return Counts{}, fmt.Errorf("open file: %w", err)
	}
	defer file.Close()

	counts := GetCounts(file)
	return counts, nil
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

func CountLines(file io.Reader) int {
	lineCount := 0

	reader := bufio.NewReader(file)

	for {
		r, _, err := reader.ReadRune()
		if err != nil {
			break
		}
		if r == '\n' {
			lineCount++
		}
	}

	return lineCount
}

func CountBytes(file io.Reader) int {
	byteCount, err := io.Copy(io.Discard, file)
	if err != nil {
		panic(err)
	}
	return int(byteCount)
}
