package server

import (
	jwtPackage "backend/pkg/jwt"
	"backend/service"
)

type Server struct {
	Service *service.Service
	JWT     jwtPackage.JWTInterface
}

type ServerInitParams struct {
	Service *service.Service
	JWT     jwtPackage.JWTInterface
}

func New(params ServerInitParams) *Server {
	return &Server{
		Service: params.Service,
		JWT:     params.JWT,
	}
}
