package auth

import "github.com/google/uuid"

type MeInput struct{}

type MeOutput struct{}

type RegisterInput struct {
	Email       string `json:"email"`
	FullName    string `json:"full_name"`
	Password    string `json:"password"`
	PhoneNumber string `json:"phone_number"`
}

type RegisterOutput struct {
	UserID uuid.UUID
}

type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginOutput struct {
	UserId   uuid.UUID
	FullName string
	Email    string
}

type VerifyInput struct {
	Email             string `json:"email"`
	VerificationToken string `json:"verificationToken"`
}

type VerifyOutput struct {
	UserID uuid.UUID
}
