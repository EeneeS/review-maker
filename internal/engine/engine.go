package engine

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type engine struct {
	hashes			 []string
	baseBranch 	 string
	targetBranch string
}

func getBaseBranch() (string, error) {
    cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
    output, err := cmd.Output()
    if err != nil {
        return "", err
    }

    branch := strings.TrimSpace(string(output))
		return branch, nil
}

func New(hashes []string) (*engine, error) {
	bb, err := getBaseBranch() 
	if err != nil {
		return nil, err
	}

	return &engine{
		hashes: hashes,
		baseBranch: bb,
	}, nil
}

func (e *engine) ProcessReview() error {
	tb, err := getTargetBranch()
	if err != nil {
		return fmt.Errorf("failed to get target branch: %w", err)
	}
	e.targetBranch = tb

	// Create the targetBranch from the baseBranch.
	cmd := exec.Command("git", "branch", e.targetBranch, e.baseBranch)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to create branch %s from %s: %w", tb, e.baseBranch, err)
	}

	// Should stash (if any) changes on the current branch before going to the targetBranch
	// Will pop them after the process is done ...

	// Command(s) to cherry pick the commits in the correct order ...
	// git cherry-pick $(git rev-list --reverse [selectedHashes])

	return nil
}

func getTargetBranch() (string, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter the target branch: ")

	input, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	
	input = strings.TrimSpace(input)
	return input, nil
}
