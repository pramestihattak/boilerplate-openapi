package server

import (
	"backend/api"
	"backend/pkg/jwt"
	"backend/service/auth"
	"encoding/json"
	"net/http"
)

// (POST /v1/auth/login)
func (s *Server) Login(w http.ResponseWriter, r *http.Request) {
	var req api.AuthLoginRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		httpError(w, ErrUnableToParseJSON)
		return
	}

	out, err := s.Service.Auth.Login(r.Context(), auth.LoginInput{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		httpError(w, err)
		return
	}

	token, err := s.JWT.Sign(jwt.Auth{
		UserID:   out.UserId.String(),
		FullName: out.FullName,
		Email:    out.Email,
	})
	if err != nil {
		httpError(w, auth.ErrFailedToLogin)
		return
	}

	// set header
	w.Header().Set("X-Jwt", token)
	json.NewEncoder(w).Encode(api.AuthLoginResponse{
		Token: token,
	})
}
