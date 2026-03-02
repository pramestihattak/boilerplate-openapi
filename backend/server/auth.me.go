package server

import (
	"backend/api"
	"backend/service/auth"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

// (GET /v1/auth/me)
func (s *Server) Me(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(UserID).(string)
	if !ok || userID == "" {
		httpError(w, ErrUnauthorized)
		return
	}

	out, err := s.Service.Auth.Me(r.Context(), auth.MeInput{
		UserID: userID,
	})
	if err != nil {
		httpError(w, err)
		return
	}

	parsedID, err := uuid.Parse(out.UserID)
	if err != nil {
		httpError(w, ErrUnauthorized)
		return
	}

	resp := api.AuthMeResponse{
		UserId:   parsedID,
		Email:    out.Email,
		FullName: out.FullName,
		Verified: out.Verified,
	}
	if out.PhoneNumber != "" {
		resp.PhoneNumber = &out.PhoneNumber
	}

	json.NewEncoder(w).Encode(resp)
}
