package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"golang.org/x/term"
)

// CommitPicker manages the interactive commit selection interface
type CommitPicker struct {
	commits     []string
	page        int
	pageSize    int
	selectedIdx int
	fd          int
	oldState    *term.State
}

// NewCommitPicker creates a new commit picker instance
func NewCommitPicker(commits []string, pageSize int) *CommitPicker {
	return &CommitPicker{
		commits:     commits,
		page:        0,
		pageSize:    pageSize,
		selectedIdx: 0,
	}
}

// Run starts the interactive picker and returns the selected commit index
func (cp *CommitPicker) Run() (int, error) {
	if err := cp.setupTerminal(); err != nil {
		return -1, err
	}
	defer cp.restoreTerminal()

	return cp.eventLoop()
}

func (cp *CommitPicker) setupTerminal() error {
	cp.fd = int(os.Stdin.Fd())
	oldState, err := term.MakeRaw(cp.fd)
	if err != nil {
		return err
	}
	cp.oldState = oldState
	return nil
}

func (cp *CommitPicker) restoreTerminal() {
	if cp.oldState != nil {
		term.Restore(cp.fd, cp.oldState)
	}
}

func (cp *CommitPicker) eventLoop() (int, error) {
	buf := make([]byte, 1)

	for {
		cp.render()

		n, err := os.Stdin.Read(buf)
		if err != nil {
			return -1, err
		}
		if n == 0 {
			continue
		}

		action := cp.handleInput(buf[0])
		if action == actionQuit {
			return -1, nil
		}
		if action == actionSelect {
			return cp.getAbsoluteIndex(), nil
		}
	}
}

type inputAction int

const (
	actionNone inputAction = iota
	actionQuit
	actionSelect
)

func (cp *CommitPicker) handleInput(key byte) inputAction {
	switch key {
	case 'q':
		return actionQuit
		case 'B': // Down arrow
		cp.moveDown()
		case 'A': // Up arrow
		cp.moveUp()
	case ' ':
		fmt.Print("handle select")
		case '\r': // Enter
		return actionSelect
	}
	return actionNone
}

func (cp *CommitPicker) moveDown() {
	visibleCommits := cp.getVisibleCommits()
	if cp.selectedIdx < len(visibleCommits)-1 {
		cp.selectedIdx++
	} else if cp.hasNextPage() {
		cp.page++
		cp.selectedIdx = 0
	}
}

func (cp *CommitPicker) moveUp() {
	if cp.selectedIdx > 0 {
		cp.selectedIdx--
	} else if cp.page > 0 {
		cp.page--
		prevPageCommits := cp.getVisibleCommits()
		cp.selectedIdx = len(prevPageCommits) - 1
	}
}

func (cp *CommitPicker) hasNextPage() bool {
	return (cp.page+1)*cp.pageSize < len(cp.commits)
}

func (cp *CommitPicker) getVisibleCommits() []string {
	return paginate(cp.commits, cp.page, cp.pageSize)
}

func (cp *CommitPicker) getAbsoluteIndex() int {
	return cp.page*cp.pageSize + cp.selectedIdx
}

func (cp *CommitPicker) render() {
	clearScreen()
	visibleCommits := cp.getVisibleCommits()

	if cp.page > 0 {
		fmt.Print("...more...\r\n")
	}

	for i, msg := range visibleCommits {
		if i == cp.selectedIdx {
			fmt.Print("> " + msg + "\r\n")
		} else {
			fmt.Print("  " + msg + "\r\n")
		}
	}

	if cp.hasNextPage() {
		fmt.Print("...more...\r\n")
	}
}

func clearScreen() {
	fmt.Print("\033[H\033[2J")
}

func paginate[T any](items []T, page int, pageSize int) []T {
	start := min(page*pageSize, len(items))
	end := min(start+pageSize, len(items))
	return items[start:end]
}

// CommitRepository handles git operations
type CommitRepository struct{}

func (cr *CommitRepository) GetCommits() ([]string, error) {
	out, err := exec.Command("git", "log", "--pretty=format:%s").Output()
	if err != nil {
		return nil, err
	}
	messages := strings.Split(string(out), "\n")
	return messages, nil
}

func (cr *CommitRepository) GetMockCommits() []string {
	return []string{
		"1", "2", "3", "4", "5", "6", "7", "8", "9", "10",
		"11", "12", "13", "14", "15", "16", "17", "18", "19", "20",
		"21", "22", "23", "24",
	}
}

func main() {
	repo := &CommitRepository{}
	// commits, err := repo.GetCommits()
	commits := repo.GetMockCommits()

	picker := NewCommitPicker(commits, 10)
	selectedIdx, err := picker.Run()

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("selected", commits[selectedIdx])
}
