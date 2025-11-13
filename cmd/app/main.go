package main

import (
	"fmt"
	"os"

	"github.com/EeneeS/review-maker/internal/picker"
	"github.com/EeneeS/review-maker/internal/repository"
)

func main() {
	repo := repository.New()
	commits, err := repo.GetCommits()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching commits: %v\n", err)
		os.Exit(1)
	}

	p := picker.New(commits, 10)
	selectedHashes, err := p.Run()

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	if len(selectedHashes) > 0 {
		fmt.Println("\nSelected commit hashes:")
		for hash := range selectedHashes {
			fmt.Printf("  - %s\n", hash)
		}
	} else {
		fmt.Println("\nSelection cancelled or no commits selected")
	}
}
