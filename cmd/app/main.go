package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"golang.org/x/term"
)

type Commit struct {
	Hash    string
	Subject string
}

type CommitPicker struct {
	commits      []Commit
	page         int
	pageSize     int
	selectedIdx  int
	fd           int
	oldState     *term.State
	selectedMap  map[string]bool // not really sure about this approach ...
}

func NewCommitPicker(commits []Commit, pageSize int) *CommitPicker {
	return &CommitPicker{
		commits:     commits,
		page:        0,
		pageSize:    pageSize,
		selectedIdx: 0,
		selectedMap: make(map[string]bool),
	}
}

func (cp *CommitPicker) Run() (map[string]bool, error) {
	if err := cp.setupTerminal(); err != nil {
		return nil, err
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

func (cp *CommitPicker) eventLoop() (map[string]bool, error) {
	buf := make([]byte, 1)

	for {
		cp.render()

		n, err := os.Stdin.Read(buf)
		if err != nil {
			return nil, err
		}
		if n == 0 {
			continue
		}

		action := cp.handleInput(buf[0])
		if action == actionQuit {
			return nil, nil
		}
		if action == actionConfirm {
			return cp.selectedMap, nil
		}
	}
}

type inputAction int

const (
	actionNone inputAction = iota
	actionQuit
	actionConfirm
)

func (cp *CommitPicker) handleInput(key byte) inputAction {
	switch key {
	case 'q':
		return actionQuit
	case 'B', 'j': // Down arrow || j
		cp.moveDown()
	case 'A', 'k': // Up arrow || k
		cp.moveUp()
	case ' ': // Space to toggle selection
		cp.toggleSelection()
	case '\r': // Enter to confirm selections
		return actionConfirm
	}
	return actionNone
}

func (cp *CommitPicker) toggleSelection() {
	absIdx := cp.getAbsoluteIndex()
	if absIdx >= 0 && absIdx < len(cp.commits) {
		hash := cp.commits[absIdx].Hash
		if cp.selectedMap[hash] {
			delete(cp.selectedMap, hash)
		} else {
			cp.selectedMap[hash] = true
		}
	}
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

func (cp *CommitPicker) getVisibleCommits() []Commit {
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
	} else {
		fmt.Print("\r\n")
	}

	for i, commit := range visibleCommits {
		isSelected := cp.selectedMap[commit.Hash]
		
		prefix := "  "
		if i == cp.selectedIdx {
			prefix = "> "
		}
		
		checkbox := "[ ]"
		if isSelected {
			checkbox = "[x]"
		}
		
		fmt.Printf("%s%s %s %s\r\n", prefix, checkbox, commit.Hash, commit.Subject)
	}

	if cp.hasNextPage() {
		fmt.Print("...more...\r\n")
	}
	
	fmt.Printf("\r\nSelected: %d commits | Space: toggle | Enter: confirm | q: quit\r\n", len(cp.selectedMap))
}

func clearScreen() {
	fmt.Print("\033[H\033[2J")
}

func paginate(items []Commit, page int, pageSize int) []Commit {
	start := min(page*pageSize, len(items))
	end := min(start+pageSize, len(items))
	return items[start:end]
}

// CommitRepository handles git operations
type CommitRepository struct{}

func (cr *CommitRepository) GetCommits() ([]Commit, error) {
	out, err := exec.Command("git", "log", "-n 50", "--pretty=format:%h\t%s").Output()
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(out), "\n")
	commits := make([]Commit, 0, len(lines))

	for _, line := range lines {
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, "\t", 2)
		if len(parts) == 2 {
			commits = append(commits, Commit{
				Hash:    parts[0],
				Subject: parts[1],
			})
		}
	}

	return commits, nil
}

func (cr *CommitRepository) GetMockCommits() []Commit { // Thanks AI ;)
	return []Commit{
		{Hash: "a1b2c3d", Subject: "1"},
		{Hash: "e4f5g6h", Subject: "2"},
		{Hash: "i7j8k9l", Subject: "3"},
		{Hash: "m0n1o2p", Subject: "4"},
		{Hash: "q3r4s5t", Subject: "5"},
		{Hash: "u6v7w8x", Subject: "6"},
		{Hash: "y9z0a1b", Subject: "7"},
		{Hash: "c2d3e4f", Subject: "8"},
		{Hash: "g5h6i7j", Subject: "9"},
		{Hash: "k8l9m0n", Subject: "10"},
		{Hash: "o1p2q3r", Subject: "11"},
		{Hash: "s4t5u6v", Subject: "12"},
		{Hash: "w7x8y9z", Subject: "13"},
		{Hash: "a0b1c2d", Subject: "14"},
		{Hash: "e3f4g5h", Subject: "15"},
		{Hash: "i6j7k8l", Subject: "16"},
		{Hash: "m9n0o1p", Subject: "17"},
		{Hash: "q2r3s4t", Subject: "18"},
		{Hash: "u5v6w7x", Subject: "19"},
		{Hash: "y8z9a0b", Subject: "20"},
		{Hash: "c1d2e3f", Subject: "21"},
		{Hash: "g4h5i6j", Subject: "22"},
		{Hash: "k7l8m9n", Subject: "23"},
		{Hash: "o0p1q2r", Subject: "24"},
	}
}

func main() {
	repo := &CommitRepository{}
	commits, err := repo.GetCommits()
	// commits := repo.GetMockCommits()

	picker := NewCommitPicker(commits, 10)
	selectedHashes, err := picker.Run()

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
