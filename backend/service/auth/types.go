package auth

import "github.com/google/uuid"

type MeInput struct {
	UserID string
}

type MeOutput struct {
	UserID      string
	Email       string
	FullName    string
	Verified    bool
	PhoneNumber string
}

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
	UserID   uuid.UUID
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
