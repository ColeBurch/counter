package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"unicode"
)

type Counts struct {
	Lines int
	Words int
	Bytes int
}

func (c *Counts) Add(other Counts) {
	c.Lines += other.Lines
	c.Words += other.Words
	c.Bytes += other.Bytes
}

func (c Counts) Print(w io.Writer, opts DisplayOptions, suffixes ...string) {
	xs := []string{}

	if opts.ShouldShowLines() {
		xs = append(xs, strconv.Itoa(c.Lines))
	}

	if opts.ShouldShowWords() {
		xs = append(xs, strconv.Itoa(c.Words))
	}

	if opts.ShouldShowBytes() {
		xs = append(xs, strconv.Itoa(c.Bytes))
	}

	xs = append(xs, suffixes...)

	line := strings.Join(xs, " ")

	fmt.Println(line)
}

func GetCounts(file io.Reader) Counts {
	res := Counts{}

	isInsideWord := false

	reader := bufio.NewReader(file)

	for {
		r, size, err := reader.ReadRune()
		if err != nil {
			break
		}
		res.Bytes += size
		if r == '\n' {
			res.Lines++
		}

		isSpace := unicode.IsSpace(r)

		if !isSpace && !isInsideWord {
			res.Words++
		}

		isInsideWord = !isSpace
	}

	return res
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
