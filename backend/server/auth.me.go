package server

import (
	"backend/api"
	"encoding/json"
	"net/http"
)

// (GET /v1/auth/me)
func (s *Server) Me(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID")
	json.NewEncoder(w).Encode(api.AuthMeResponse{
		Message: userID.(string),
	})
}
