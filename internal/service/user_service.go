package service

import (
	"context"
	"pull_requests_service/internal/domain"
)

type UserService struct {
	userRepo domain.UserRepository
	prRepo   domain.PRRepository
}

func NewUserService(userRepo domain.UserRepository, prRepo domain.PRRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
		prRepo:   prRepo,
	}
}

func (s *UserService) SetUserActive(ctx context.Context, userID string, isActive bool) (*domain.TeamMember, string, error) {
	user, err := s.userRepo.GetUser(userID)
	if err != nil {
		return nil, "", err
	}

	teamName, err := s.userRepo.GetUserTeam(userID)
	if err != nil {
		return nil, "", err
	}

	if err := s.userRepo.SetUserActive(userID, isActive); err != nil {
		return nil, "", err
	}

	user.IsActive = isActive
	return user, teamName, nil
}

func (s *UserService) GetUserReviews(ctx context.Context, userID string) ([]*domain.PullRequest, error) {
	if _, err := s.userRepo.GetUser(userID); err != nil {
		return nil, err
	}

	prs, err := s.prRepo.GetPRsByReviewer(userID)
	if err != nil {
		return nil, err
	}

	return prs, nil
}
