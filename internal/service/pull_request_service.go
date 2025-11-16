package service

import (
	"context"
	"math/rand"
	"pull_requests_service/internal/domain"
	"pull_requests_service/internal/dto"
	"slices"
	"time"
)

type PRService struct {
	prRepo   domain.PRRepository
	userRepo domain.UserRepository
}

func NewPRService(prRepo domain.PRRepository, userRepo domain.UserRepository) *PRService {
	return &PRService{
		prRepo:   prRepo,
		userRepo: userRepo,
	}
}

func (s *PRService) CreatePR(ctx context.Context, req dto.CreatePRRequest) (*domain.PullRequest, error) {
	exists, err := s.prRepo.PRExists(req.PullRequestID)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, domain.ErrPRExists
	}

	teamName, err := s.userRepo.GetUserTeam(req.AuthorID)
	if err != nil {
		return nil, err
	}

	members, err := s.userRepo.GetActiveTeamMembers(teamName)
	if err != nil {
		return nil, err
	}

	var reviewers []string
	candidates := make([]domain.TeamMember, 0, len(members))

	for _, member := range members {
		if member.UserID != req.AuthorID {
			candidates = append(candidates, member)
		}
	}

	rand.New(rand.NewSource(time.Now().UnixNano()))
	rand.Shuffle(len(candidates), func(i, j int) {
		candidates[i], candidates[j] = candidates[j], candidates[i]
	})

	for i := 0; i < len(candidates) && i < 2; i++ {
		reviewers = append(reviewers, candidates[i].UserID)
	}

	pr := &domain.PullRequest{
		PullRequestID:     req.PullRequestID,
		PullRequestName:   req.PullRequestName,
		AuthorID:          req.AuthorID,
		Status:            domain.PRStatusOpen,
		AssignedReviewers: reviewers,
	}

	if err := s.prRepo.CreatePR(pr); err != nil {
		return nil, err
	}

	return pr, nil
}

func (s *PRService) MergePR(ctx context.Context, prID string) (*domain.PullRequest, error) {
	pr, err := s.prRepo.GetPR(prID)
	if err != nil {
		return nil, err
	}

	if pr.Status == domain.PRStatusMerged {
		return pr, nil
	}

	pr.Status = domain.PRStatusMerged
	now := time.Now()
	pr.MergedAt = &now

	if err := s.prRepo.UpdatePR(pr); err != nil {
		return nil, err
	}

	return pr, nil
}

func (s *PRService) ReassignPR(ctx context.Context, prID, oldUserID string) (string, *domain.PullRequest, error) {
	pr, err := s.prRepo.GetPR(prID)
	if err != nil {
		return "", nil, err
	}

	if pr.Status == domain.PRStatusMerged {
		return "", nil, domain.ErrPRMerged
	}

	if !contains(pr.AssignedReviewers, oldUserID) {
		return "", nil, domain.ErrNotAssigned
	}

	teamName, err := s.userRepo.GetUserTeam(oldUserID)
	if err != nil {
		return "", nil, err
	}

	members, err := s.userRepo.GetActiveTeamMembers(teamName)
	if err != nil {
		return "", nil, err
	}

	var candidates []string
	for _, member := range members {
		if member.UserID != oldUserID && member.UserID != pr.AuthorID && !contains(pr.AssignedReviewers, member.UserID) {
			candidates = append(candidates, member.UserID)
		}
	}

	if len(candidates) == 0 {
		return "", nil, domain.ErrNoCandidate
	}

	rand.New(rand.NewSource(time.Now().UnixNano()))
	newUserID := candidates[rand.Intn(len(candidates))]

	for i, reviewer := range pr.AssignedReviewers {
		if reviewer == oldUserID {
			pr.AssignedReviewers[i] = newUserID
			break
		}
	}

	if err := s.prRepo.UpdatePR(pr); err != nil {
		return "", nil, err
	}

	return newUserID, pr, nil
}

func contains(slice []string, item string) bool {
	return slices.Contains(slice, item)
}
