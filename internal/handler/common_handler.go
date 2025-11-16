package handler

import (
	"encoding/json"
	"net/http"
	"pull_requests_service/internal/dto"
)

func writeError(w http.ResponseWriter, status int, code, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(dto.NewErrorResponse(code, message))
}
