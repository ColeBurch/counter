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
		totals.Add(counts)
		counts.Print(os.Stdout, filename)
	}

	if len(filenames) == 0 {
		counts := GetCounts(os.Stdin)
		counts.Print(os.Stdout)
	}

	if len(filenames) > 0 {
		totals.Print(os.Stdout, "total")
	}

	if didError {
		os.Exit(1)
	}
}
