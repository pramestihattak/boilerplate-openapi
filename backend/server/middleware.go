package server

import (
	"context"
	"net/http"
)

const (
	UserID = "userID"
)

var restrictedAuthEndpoints = map[string]bool{
	"/v1/auth/me": true,
}

func (s *Server) WithAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !restrictedAuthEndpoints[r.RequestURI] {
			next.ServeHTTP(w, r)
			return
		}

		token := r.Header.Get("Authorization")
		if token == "" {
			httpError(w, ErrUnauthorized)
			return
		}

		if !s.JWT.IsValidToken(token) {
			httpError(w, ErrUnauthorized)
			return
		}

		auth, err := s.JWT.GetClaims(token)
		if err != nil {
			httpError(w, ErrUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), UserID, auth.UserID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
