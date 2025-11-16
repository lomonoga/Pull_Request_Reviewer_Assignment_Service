package service

import (
	"context"
	"pull_requests_service/internal/domain"
	"pull_requests_service/internal/dto"
)

type TeamService struct {
	teamRepo domain.TeamRepository
	userRepo domain.UserRepository
}

func NewTeamService(teamRepo domain.TeamRepository, userRepo domain.UserRepository) *TeamService {
	return &TeamService{
		teamRepo: teamRepo,
		userRepo: userRepo,
	}
}

func (s *TeamService) AddTeam(ctx context.Context, req dto.AddTeamRequest) (*domain.Team, error) {
	exists, err := s.teamRepo.TeamExists(req.TeamName)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, domain.ErrTeamExists
	}

	team := &domain.Team{
		TeamName: req.TeamName,
		Members:  make([]domain.TeamMember, len(req.Members)),
	}

	for i, member := range req.Members {
		team.Members[i] = domain.TeamMember{
			UserID:   member.UserID,
			Username: member.Username,
			IsActive: member.IsActive,
		}

		if err := s.userRepo.CreateOrUpdateUser(&team.Members[i]); err != nil {
			return nil, err
		}
	}

	if err := s.teamRepo.CreateTeam(team); err != nil {
		return nil, err
	}

	return team, nil
}

func (s *TeamService) GetTeam(ctx context.Context, teamName string) (*domain.Team, error) {
	team, err := s.teamRepo.GetTeam(teamName)
	if err != nil {
		return nil, err
	}
	return team, nil
}
