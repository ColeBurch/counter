package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	log.SetFlags(0)

	filenames := os.Args[1:]
	totals := Counts{}
	didError := false

	for _, filename := range filenames {
		counts, err := CountFile(filename)
		if err != nil {
			didError = true
			fmt.Printf("%s: %s\n", filename, err)
			continue
		}
		totals = Counts{
			Lines: totals.Lines + counts.Lines,
			Words: totals.Words + counts.Words,
			Bytes: totals.Bytes + counts.Bytes,
		}
		fmt.Printf("%s: %d %d %d\n", filename, counts.Lines, counts.Words, counts.Bytes)
	}

	if len(filenames) == 0 {
		counts := GetCounts(os.Stdin)
		fmt.Println(counts.Lines, counts.Words, counts.Bytes)
	}

	if len(filenames) > 0 {
		fmt.Printf("Total count: %d %d %d\n", totals.Lines, totals.Words, totals.Bytes)
	}

	if didError {
		os.Exit(1)
	}
}
