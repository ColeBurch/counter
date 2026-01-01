package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"sync"
	"unicode"
)

type Counts struct {
	Lines int
	Words int
	Bytes int
}

type FileCountsResult struct {
	Filename string
	Counts   Counts
}

func (c *Counts) Add(other Counts) {
	c.Lines += other.Lines
	c.Words += other.Words
	c.Bytes += other.Bytes
}

func (c Counts) Print(w io.Writer, opts DisplayOptions, suffixes ...string) {
	stats := []string{}

	if opts.Headers {
		if opts.ShouldShowLines() {
			stats = append(stats, "Lines")
		}
		if opts.ShouldShowWords() {
			stats = append(stats, "Words")
		}
		if opts.ShouldShowBytes() {
			stats = append(stats, "Bytes")
		}
	} else {
		if opts.ShouldShowLines() {
			stats = append(stats, strconv.Itoa(c.Lines))
		}

		if opts.ShouldShowWords() {
			stats = append(stats, strconv.Itoa(c.Words))
		}

		if opts.ShouldShowBytes() {
			stats = append(stats, strconv.Itoa(c.Bytes))
		}
	}

	line := strings.Join(stats, "\t") + "\t"

	fmt.Fprint(w, line)

	suffixStr := strings.Join(suffixes, "\t")
	if suffixStr != "" {
		fmt.Fprintf(w, " %s", suffixStr)
	}

	fmt.Fprint(w, "\n")
}

func GetCounts(file io.Reader) Counts {
	bytesReader, bytesWriter := io.Pipe()
	wordsReader, wordsWriter := io.Pipe()
	linesReader, linesWriter := io.Pipe()

	w := io.MultiWriter(bytesWriter, wordsWriter, linesWriter)

	chBytes := make(chan int)
	chWords := make(chan int)
	chLines := make(chan int)

	go func() {
		defer close(chBytes)
		chBytes <- CountBytes(bytesReader)
	}()

	go func() {
		defer close(chWords)
		chWords <- CountWords(wordsReader)
	}()

	go func() {
		defer close(chLines)
		chLines <- CountLines(linesReader)
	}()

	io.Copy(w, file)
	bytesWriter.Close()
	wordsWriter.Close()
	linesWriter.Close()

	return Counts{
		Bytes: <-chBytes,
		Words: <-chWords,
		Lines: <-chLines,
	}
}

func GetCountsSinglePass(file io.Reader) Counts {
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

func CountFiles(filenames []string) (<-chan FileCountsResult, <-chan error) {
	ch := make(chan FileCountsResult)
	errCh := make(chan error)

	wg := sync.WaitGroup{}
	wg.Add(len(filenames))

	for _, filename := range filenames {
		go func() {
			defer wg.Done()

			file, err := os.Open(filename)
			if err != nil {
				errCh <- err
				return
			}
			defer file.Close()

			counts := GetCountsSinglePass(file)
			ch <- FileCountsResult{filename, counts}
		}()
	}

	go func() {
		wg.Wait()
		close(ch)
		close(errCh)
	}()

	return ch, errCh
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
