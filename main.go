package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

type DisplayOptions struct {
	ShowLines bool
	ShowWords bool
	ShowBytes bool
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
	}

	flag.BoolVar(&opts.ShowLines, "l", false, "show line count")
	flag.BoolVar(&opts.ShowWords, "w", false, "show word count")
	flag.BoolVar(&opts.ShowBytes, "c", false, "show byte count")

	flag.Parse()

	log.SetFlags(0)

	filenames := flag.Args()
	totals := Counts{}
	didError := false

	for _, filename := range filenames {
		counts, err := CountFile(filename)
		if err != nil {
			didError = true
			fmt.Printf("%s: %s\n", filename, err)
			continue
		}
		totals.Add(counts)
		counts.Print(os.Stdout, opts, filename)
	}

	if len(filenames) == 0 {
		counts := GetCounts(os.Stdin)
		counts.Print(os.Stdout, opts)
	}

	if len(filenames) > 0 {
		totals.Print(os.Stdout, opts, "total")
	}

	if didError {
		os.Exit(1)
	}
}
