package picker

import "fmt"

func (p *Picker) render() {
	clearScreen()
	visibleCommits := p.getVisibleCommits()

	if p.page > 0 {
		fmt.Print("...more...\r\n")
	} else {
		fmt.Print("\r\n")
	}

	for i, commit := range visibleCommits {
		isSelected := p.selectedMap[commit.Hash]

		prefix := "  "
		if i == p.selectedIdx {
			prefix = "> "
		}

		checkbox := "[ ]"
		if isSelected {
			checkbox = "[x]"
		}

		fmt.Printf("%s%s %s %s\r\n", prefix, checkbox, commit.Hash, commit.Subject)
	}

	if p.hasNextPage() {
		fmt.Print("...more...\r\n")
	}

	fmt.Printf("\r\nSelected: %d commits | Space: toggle | Enter: confirm | q: quit\r\n", len(p.selectedMap))
}

func clearScreen() {
	fmt.Print("\033[H\033[2J")
}
