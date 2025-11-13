package entity

import "time"

type PullRequest struct {
	ID              string
	Name            string
	AuthorID        string
	Status          PRStatus
	AssignReviewers []User
	CreatedAt       *time.Time
	MergedAt        *time.Time
}

type PullRequestShort struct {
	ID       string
	Name     string
	AuthorID string
	Status   PRStatus
}

type PRStatus string

const (
	OPEN   PRStatus = "OPEN"
	MERGED PRStatus = "MERGED"
)
