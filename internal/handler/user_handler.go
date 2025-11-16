package handler

import (
	"encoding/json"
	"net/http"
	"pull_requests_service/internal/domain"
	"pull_requests_service/internal/dto"
	"pull_requests_service/internal/service"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) SetUserActive(w http.ResponseWriter, r *http.Request) {
	var req dto.SetUserActiveRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "NOT_FOUND", "resource not found")
		return
	}

	user, teamName, err := h.userService.SetUserActive(r.Context(), req.UserID, req.IsActive)
	if err != nil {
		switch err {
		case domain.ErrUserNotFound:
			writeError(w, http.StatusNotFound, "NOT_FOUND", "user not found")
		default:
			writeError(w, http.StatusInternalServerError, "NOT_FOUND", err.Error())
		}
		return
	}

	response := dto.SetUserActiveResponse{
		User: dto.User{
			UserID:   user.UserID,
			Username: user.Username,
			TeamName: teamName,
			IsActive: user.IsActive,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *UserHandler) GetUserReviews(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		writeError(w, http.StatusBadRequest, "NOT_FOUND", "user_id is required")
		return
	}

	prs, err := h.userService.GetUserReviews(r.Context(), userID)
	if err != nil {
		switch err {
		case domain.ErrUserNotFound:
			writeError(w, http.StatusNotFound, "NOT_FOUND", "resource not found")
		default:
			writeError(w, http.StatusInternalServerError, "NOT_FOUND", err.Error())
		}
		return
	}

	response := dto.GetUserReviewsResponse{
		UserID:       userID,
		PullRequests: make([]dto.PullRequestShort, len(prs)),
	}

	for i, pr := range prs {
		response.PullRequests[i] = dto.PullRequestShort{
			PullRequestID:   pr.PullRequestID,
			PullRequestName: pr.PullRequestName,
			AuthorID:        pr.AuthorID,
			Status:          string(pr.Status),
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
