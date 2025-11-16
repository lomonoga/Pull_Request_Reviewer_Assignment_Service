package domain

import "errors"

var (
	ErrTeamExists   = errors.New("team already exists")
	ErrTeamNotFound = errors.New("team not found")
	ErrUserNotFound = errors.New("user not found")
	ErrPRExists     = errors.New("PR already exists")
	ErrPRNotFound   = errors.New("PR not found")
	ErrPRMerged     = errors.New("PR is merged")
	ErrNotAssigned  = errors.New("cannot reassign on merged PR")
	ErrNoCandidate  = errors.New("no active replacement candidate in team")
)
