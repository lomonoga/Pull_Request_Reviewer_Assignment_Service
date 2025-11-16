package handler

import (
	"encoding/json"
	"net/http"
	"pull_requests_service/internal/domain"
	"pull_requests_service/internal/dto"
	"pull_requests_service/internal/service"
)

type TeamHandler struct {
	teamService *service.TeamService
}

func NewTeamHandler(teamService *service.TeamService) *TeamHandler {
	return &TeamHandler{teamService: teamService}
}

func (h *TeamHandler) AddTeam(w http.ResponseWriter, r *http.Request) {
	var req dto.AddTeamRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "NOT_FOUND", "invalid request body")
		return
	}

	team, err := h.teamService.AddTeam(r.Context(), req)
	if err != nil {
		switch err {
		case domain.ErrTeamExists:
			writeError(w, http.StatusBadRequest, "TEAM_EXISTS", "team_name already exists")
		default:
			writeError(w, http.StatusInternalServerError, "NOT_FOUND", err.Error())
		}
		return
	}

	response := dto.AddTeamResponse{
		Team: dto.Team{
			TeamName: team.TeamName,
			Members:  make([]dto.TeamMember, len(team.Members)),
		},
	}

	for i, member := range team.Members {
		response.Team.Members[i] = dto.TeamMember{
			UserID:   member.UserID,
			Username: member.Username,
			IsActive: member.IsActive,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func (h *TeamHandler) GetTeam(w http.ResponseWriter, r *http.Request) {
	teamName := r.URL.Query().Get("team_name")
	if teamName == "" {
		writeError(w, http.StatusBadRequest, "NOT_FOUND", "team_name is required")
		return
	}

	team, err := h.teamService.GetTeam(r.Context(), teamName)
	if err != nil {
		switch err {
		case domain.ErrTeamNotFound:
			writeError(w, http.StatusNotFound, "NOT_FOUND", "resource not found")
		default:
			writeError(w, http.StatusInternalServerError, "NOT_FOUND", err.Error())
		}
		return
	}

	response := dto.GetTeamResponse{
		TeamName: team.TeamName,
		Members:  make([]dto.TeamMember, len(team.Members)),
	}

	for i, member := range team.Members {
		response.Members[i] = dto.TeamMember{
			UserID:   member.UserID,
			Username: member.Username,
			IsActive: member.IsActive,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
