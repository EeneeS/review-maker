package repository

import (
	"os/exec"
	"strings"

	"github.com/EeneeS/review-maker/internal/models"
)

type GitRepository struct{}


func (r *GitRepository) GetCommits() ([]models.Commit, error) {
	out, err := exec.Command("git", "log", "-n 50", "--pretty=format:%h\t%s").Output()
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(out), "\n")
	commits := make([]models.Commit, 0, len(lines))

	for _, line := range lines {
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, "\t", 2)
		if len(parts) == 2 {
			commits = append(commits, models.Commit{
				Hash:    parts[0],
				Subject: parts[1],
			})
		}
	}

	return commits, nil
}


