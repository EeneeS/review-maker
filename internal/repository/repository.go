package repository

import "github.com/EeneeS/review-maker/internal/models"

type Repository interface {
	GetCommits() ([]models.Commit, error)
}

func New() *GitRepository {
	return &GitRepository{}
}

func NewMock() *MockRepository {
	return &MockRepository{}
}
