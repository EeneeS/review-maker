package repository

import "github.com/EeneeS/review-maker/internal/models"

type MockRepository struct {}

func (r *MockRepository) GetCommits() ([]models.Commit, error) {
	return []models.Commit{
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
	}, nil
}


