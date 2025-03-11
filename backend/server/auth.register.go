package server

import (
	"backend/api"
	"backend/service/auth"
	"encoding/json"
	"fmt"
	"net/http"
)

// (POST /v1/auth/register)
func (s *Server) Register(w http.ResponseWriter, r *http.Request) {
	var req api.AuthRegisterRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		httpError(w, ErrUnableToParseJSON)
		return
	}

	out, err := s.Service.Auth.Register(r.Context(), auth.RegisterInput{
		Email:       req.Email,
		FullName:    req.FullName,
		Password:    req.Password,
		PhoneNumber: req.PhoneNumber,
	})
	if err != nil {
		httpError(w, err)
		return
	}

	json.NewEncoder(w).Encode(api.AuthRegisterResponse{
		Message: fmt.Sprintf("success: %v", out.UserID),
	})
}
