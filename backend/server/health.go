package server

import (
	"backend/api"
	"encoding/json"
	"net/http"
)

// (GET /health)
func (s *Server) HealthCheck(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(api.HealthResponse{
		Status: "ok",
	})
}
