package picker

type inputAction int

const (
	actionNone inputAction = iota
	actionQuit
	actionConfirm
)

func (p *Picker) handleInput(key byte) inputAction {
	switch key {
	case 'q':
		return actionQuit
	case 'B', 'j': // Down arrow || j
		p.moveDown()
	case 'A', 'k': // Up arrow || k
		p.moveUp()
	case ' ': // Space to toggle selection
		p.toggleSelection()
	case '\r': // Enter to confirm selections
		return actionConfirm
	}
	return actionNone
}
