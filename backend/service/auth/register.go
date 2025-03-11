package auth

import (
	storageAuth "backend/storage/auth"
	"backend/util"
	"context"

	"github.com/google/uuid"
)

func (s *Auth) Register(ctx context.Context, input RegisterInput) (RegisterOutput, error) {
	logger := s.Logger.WithField("handler", "Auth.Register")

	exist, err := s.Storage.UserExist(ctx, input.Email)
	if err != nil {
		logger.Errorf("fail to register user: %v", err.Error())
		return RegisterOutput{}, ErrFailedToRegister
	}

	if exist > 0 {
		logger.Errorf("fail to register user: %v", "account exist")
		return RegisterOutput{}, ErrAccountAlreadyExist
	}

	hashedPassword, err := util.HashAndSalt(input.Password)
	if err != nil {
		logger.Errorf("fail to register user: %v", err.Error())
		return RegisterOutput{}, ErrFailedToRegister
	}

	userID, err := s.Storage.Register(ctx, storageAuth.Register{
		FullName:          input.FullName,
		Email:             input.Email,
		Password:          hashedPassword,
		VerificationToken: util.RandomStringGenerator(10),
		PhoneNumber:       input.PhoneNumber,
	})
	if err != nil {
		logger.Errorf("fail to register user: %v", err.Error())
		return RegisterOutput{}, ErrFailedToRegister
	}

	return RegisterOutput{
		UserID: uuid.MustParse(userID),
	}, nil
}
