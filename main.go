package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"text/tabwriter"
)

type DisplayOptions struct {
	ShowLines bool
	ShowWords bool
	ShowBytes bool
	Headers   bool
}

func (opts DisplayOptions) ShouldShowLines() bool {
	if !opts.ShowLines && !opts.ShowWords && !opts.ShowBytes {
		return true
	}

	if opts.ShowLines {
		return true
	}

	return false
}

func (opts DisplayOptions) ShouldShowWords() bool {
	if !opts.ShowLines && !opts.ShowWords && !opts.ShowBytes {
		return true
	}

	if opts.ShowWords {
		return true
	}

	return false
}

func (opts DisplayOptions) ShouldShowBytes() bool {
	if !opts.ShowLines && !opts.ShowWords && !opts.ShowBytes {
		return true
	}

	if opts.ShowBytes {
		return true
	}

	return false
}

func main() {
	opts := DisplayOptions{
		ShowLines: false,
		ShowWords: false,
		ShowBytes: false,
		Headers:   false,
	}

	flag.BoolVar(&opts.ShowLines, "l", false, "show line count")
	flag.BoolVar(&opts.ShowWords, "w", false, "show word count")
	flag.BoolVar(&opts.ShowBytes, "c", false, "show byte count")
	flag.BoolVar(&opts.Headers, "headers", false, "show headers")

	flag.Parse()

	wr := tabwriter.NewWriter(os.Stdout, 0, 8, 1, ' ', tabwriter.AlignRight)

	log.SetFlags(0)

	filenames := flag.Args()
	totals := Counts{}
	didError := false

	if opts.Headers {
		count := Counts{}
		count.Print(wr, opts)
		opts.Headers = false
	}

	ch, errCh := CountFiles(filenames)

	results := make([]FileCountsResult, len(filenames))

	filenameIndex := make(map[string]int, len(filenames))

	for i, filename := range filenames {
		filenameIndex[filename] = i
	}

	for {
		select {
		case res, open := <-ch:
			if !open {
				ch = nil
				break
			}
			idx, ok := filenameIndex[res.Filename], true
			if !ok {
				continue
			}
			results[idx] = res
		case err, open := <-errCh:
			if !open {
				errCh = nil
				break
			}
			didError = true
			fmt.Printf("%s\n", err)
		}

		if ch == nil || errCh == nil {
			break
		}
	}

	for _, result := range results {
		if result.Filename == "" {
			continue
		}
		totals.Add(result.Counts)
		result.Counts.Print(wr, opts, result.Filename)
	}

	if len(filenames) == 0 {
		counts := GetCounts(os.Stdin)
		counts.Print(wr, opts)
	}

	if len(filenames) > 0 {
		totals.Print(wr, opts, "total")
	}

	wr.Flush()

	if didError {
		os.Exit(1)
	}
}
