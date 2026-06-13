package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/devanshu0x/zentro/internal/mp4"
)

func main() {
	verbose := flag.Bool(
		"verbose",
		false,
		"show parsing details",
	)
	flag.Parse()
	if flag.NArg() != 1 {
		fmt.Println("Usage: zentro [--verbose] <path to mp4 file>")
		return
	}
	filePath := flag.Arg(0)

	info, err := os.Stat(filePath)
	if err != nil {
		fmt.Printf("No such file exists at given path: %s\n", filePath)
		return
	}
	if info.IsDir() {
		fmt.Printf("Its a directory path not a file: %s", filePath)
		return
	}

	movie, err := mp4.ParseFile(filePath, *verbose)
	if err != nil {
		log.Fatalf("Failed to parse file: %v", err)
	}

	fmt.Printf("Movie Duration: %.2f sec\n", movie.Duration)

	for _, track := range movie.Tracks {
		fmt.Printf(
			"Track %d [%s] Duration=%.2f Width=%.0f Height=%.0f\n",
			track.ID,
			track.Type,
			track.Duration,
			track.Width,
			track.Height,
		)
	}
}
