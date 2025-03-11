package auth

import "errors"

var (
	ErrAccountNotFound     = errors.New("account not found")
	ErrAccountNotVerified  = errors.New("account hasn't been verified yet")
	ErrWrongPassword       = errors.New("wrong password")
	ErrFailedToLogin       = errors.New("failed to login")
	ErrFailedToVerify      = errors.New("failed to verify")
	ErrFailedToRegister    = errors.New("failed to register")
	ErrAccountAlreadyExist = errors.New("account already exist")
)
