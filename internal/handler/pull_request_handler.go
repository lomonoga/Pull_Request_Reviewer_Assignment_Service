package handler

import (
	"encoding/json"
	"net/http"
	"pull_requests_service/internal/domain"
	"pull_requests_service/internal/dto"
	"pull_requests_service/internal/service"
	"time"
)

type PRHandler struct {
	prService *service.PRService
}

func NewPRHandler(prService *service.PRService) *PRHandler {
	return &PRHandler{prService: prService}
}

func (h *PRHandler) CreatePR(w http.ResponseWriter, r *http.Request) {
	var req dto.CreatePRRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "NOT_FOUND", "invalid request body")
		return
	}

	pr, err := h.prService.CreatePR(r.Context(), req)
	if err != nil {
		switch err {
		case domain.ErrPRExists:
			writeError(w, http.StatusConflict, "PR_EXISTS", "PR id already exists")
		case domain.ErrUserNotFound:
			writeError(w, http.StatusNotFound, "NOT_FOUND", "resource not found")
		default:
			writeError(w, http.StatusInternalServerError, "NOT_FOUND", err.Error())
		}
		return
	}

	response := h.domainPRToDTO(pr)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(dto.CreatePRResponse{PR: response})
}

func (h *PRHandler) MergePR(w http.ResponseWriter, r *http.Request) {
	var req dto.MergePRRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "NOT_FOUND", "invalid request body")
		return
	}

	pr, err := h.prService.MergePR(r.Context(), req.PullRequestID)
	if err != nil {
		switch err {
		case domain.ErrPRNotFound:
			writeError(w, http.StatusNotFound, "NOT_FOUND", "resource not found")
		default:
			writeError(w, http.StatusInternalServerError, "NOT_FOUND", err.Error())
		}
		return
	}

	response := h.domainPRToDTO(pr)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dto.MergePRResponse{PR: response})
}

func (h *PRHandler) ReassignPR(w http.ResponseWriter, r *http.Request) {
	var req dto.ReassignPRRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "NOT_FOUND", "invalid request body")
		return
	}

	newUserID, pr, err := h.prService.ReassignPR(r.Context(), req.PullRequestID, req.OldUserID)
	if err != nil {
		switch err {
		case domain.ErrPRNotFound, domain.ErrUserNotFound:
			writeError(w, http.StatusNotFound, "NOT_FOUND", "resource not found")
		case domain.ErrPRMerged:
			writeError(w, http.StatusConflict, "PR_MERGED", "cannot reassign on merged PR")
		case domain.ErrNotAssigned:
			writeError(w, http.StatusConflict, "NOT_ASSIGNED", "cannot reassign on merged PR")
		case domain.ErrNoCandidate:
			writeError(w, http.StatusConflict, "NO_CANDIDATE", "cannot reassign on merged PR")
		default:
			writeError(w, http.StatusInternalServerError, "NOT_FOUND", err.Error())
		}
		return
	}

	response := h.domainPRToDTO(pr)
	reassignResponse := dto.ReassignPRResponse{
		PR:         response,
		ReplacedBy: newUserID,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reassignResponse)
}

func (h *PRHandler) domainPRToDTO(pr *domain.PullRequest) dto.PullRequest {
	var mergedAt *string
	if pr.MergedAt != nil {
		formatted := pr.MergedAt.Format(time.RFC3339)
		mergedAt = &formatted
	}

	return dto.PullRequest{
		PullRequestID:     pr.PullRequestID,
		PullRequestName:   pr.PullRequestName,
		AuthorID:          pr.AuthorID,
		Status:            string(pr.Status),
		AssignedReviewers: pr.AssignedReviewers,
		MergedAt:          mergedAt,
	}
}
