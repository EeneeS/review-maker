package main

import (
	"fmt"
	"os"

	"github.com/EeneeS/review-maker/internal/engine"
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

	e, err := engine.New(selectedHashes)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	e.ProcessReview()

	// this is to make sure that the commits will be cherry picked in the correct order
	// git cherry-pick $(git rev-list --reverse [selectedHashes])

	// if len(selectedHashes) > 0 {
	// 	fmt.Println("\nSelected commit hashes:")
	// 	for hash := range selectedHashes {
	// 		fmt.Printf("  - %s\n", hash)
	// 	}
	// } else {
	// 	fmt.Println("\nSelection cancelled or no commits selected")
	// }
}
