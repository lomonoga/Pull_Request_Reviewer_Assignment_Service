package domain

import "time"

type PRStatus string

const (
	PRStatusOpen   PRStatus = "OPEN"
	PRStatusMerged PRStatus = "MERGED"
)

type PullRequest struct {
	PullRequestID     string
	PullRequestName   string
	AuthorID          string
	Status            PRStatus
	AssignedReviewers []string
	MergedAt          *time.Time
}

type PRRepository interface {
	CreatePR(pr *PullRequest) error
	GetPR(prID string) (*PullRequest, error)
	UpdatePR(pr *PullRequest) error
	GetPRsByReviewer(userID string) ([]*PullRequest, error)
	PRExists(prID string) (bool, error)
}
