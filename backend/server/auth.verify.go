package server

import (
	"backend/api"
	"backend/service/auth"
	"encoding/json"
	"net/http"
)

// (GET /v1/auth/verification)
func (s *Server) Verification(w http.ResponseWriter, r *http.Request, params api.VerificationParams) {
	if string(params.Email) == "" || string(params.VerificationToken) == "" {
		httpError(w, ErrVerifyMissingParams)
		return
	}
	_, err := s.Service.Auth.Verify(r.Context(), auth.VerifyInput{
		Email:             string(params.Email),
		VerificationToken: string(params.VerificationToken),
	})
	if err != nil {
		httpError(w, err)
		return
	}

	json.NewEncoder(w).Encode(api.AuthVerificationResponse{
		Message: "success",
	})
}
