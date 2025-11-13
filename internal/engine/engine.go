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

func (e *engine) Process() error {
	tb, err := getTargetBranch()
	if err != nil {
		return err
	}

	fmt.Println(tb)

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
