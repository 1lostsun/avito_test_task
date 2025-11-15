package entity

import "time"

type PullRequest struct {
	ID              string
	Name            string
	AuthorID        string
	Status          string
	AssignReviewers []string
	CreatedAt       *time.Time
	MergedAt        *time.Time
}

type PullRequestShort struct {
	ID       string
	Name     string
	AuthorID string
	Status   string
}

const (
	OPEN   = "OPEN"
	MERGED = "MERGED"
)
