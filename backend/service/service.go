package service

import (
	"backend/service/auth"
)

type Service struct {
	Auth auth.AuthService
}

type ServiceInitParams struct {
	Auth auth.AuthService
}

func New(params ServiceInitParams) *Service {
	return &Service{
		Auth: params.Auth,
	}
}
