package jwt

import "errors"

var (
	ErrInvalidToken      = errors.New("invalid token")
	ErrTokenExpired      = errors.New("token expired")
	ErrTokenNotFound     = errors.New("token not found")
	ErrFailToSignedToken = errors.New("fail to sign token")
)
