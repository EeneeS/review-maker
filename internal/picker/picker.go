package picker

import (
	"os"
	"github.com/EeneeS/review-maker/internal/models"
	"golang.org/x/term"
)

type Picker struct {
	commits     []models.Commit
	page        int
	pageSize    int
	selectedIdx int
	fd          int
	oldState    *term.State
	selectedMap map[string]bool
}

func New(commits []models.Commit, pageSize int) *Picker {
	return &Picker{
		commits:     commits,
		page:        0,
		pageSize:    pageSize,
		selectedIdx: 0,
		selectedMap: make(map[string]bool),
	}
}

func (p *Picker) Run() (map[string]bool, error) {
	if err := p.setupTerminal(); err != nil {
		return nil, err
	}
	defer p.restoreTerminal()

	return p.eventLoop()
}

func (p *Picker) setupTerminal() error {
	p.fd = int(os.Stdin.Fd())
	oldState, err := term.MakeRaw(p.fd)
	if err != nil {
		return err
	}
	p.oldState = oldState
	return nil
}

func (p *Picker) restoreTerminal() {
	if p.oldState != nil {
		term.Restore(p.fd, p.oldState)
	}
}

func (p *Picker) eventLoop() (map[string]bool, error) {
	buf := make([]byte, 1)

	for {
		p.render()

		n, err := os.Stdin.Read(buf)
		if err != nil {
			return nil, err
		}
		if n == 0 {
			continue
		}

		action := p.handleInput(buf[0])
		if action == actionQuit {
			return nil, nil
		}
		if action == actionConfirm {
			return p.selectedMap, nil
		}
	}
}

func (p *Picker) toggleSelection() {
	absIdx := p.getAbsoluteIndex()
	if absIdx >= 0 && absIdx < len(p.commits) {
		hash := p.commits[absIdx].Hash
		if p.selectedMap[hash] {
			delete(p.selectedMap, hash)
		} else {
			p.selectedMap[hash] = true
		}
	}
}

func (p *Picker) moveDown() {
	visibleCommits := p.getVisibleCommits()
	if p.selectedIdx < len(visibleCommits)-1 {
		p.selectedIdx++
	} else if p.hasNextPage() {
		p.page++
		p.selectedIdx = 0
	}
}

func (p *Picker) moveUp() {
	if p.selectedIdx > 0 {
		p.selectedIdx--
	} else if p.page > 0 {
		p.page--
		prevPageCommits := p.getVisibleCommits()
		p.selectedIdx = len(prevPageCommits) - 1
	}
}

func (p *Picker) hasNextPage() bool {
	return (p.page+1)*p.pageSize < len(p.commits)
}

func (p *Picker) getVisibleCommits() []models.Commit {
	return paginate(p.commits, p.page, p.pageSize)
}

func (p *Picker) getAbsoluteIndex() int {
	return p.page*p.pageSize + p.selectedIdx
}

func paginate(items []models.Commit, page int, pageSize int) []models.Commit {
	start := min(page*pageSize, len(items))
	end := min(start+pageSize, len(items))
	return items[start:end]
}
