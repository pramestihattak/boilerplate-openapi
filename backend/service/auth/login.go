package auth

import (
	"backend/util"
	"context"

	storageAuth "backend/storage/auth"

	"github.com/google/uuid"
)

func (s *Auth) Login(ctx context.Context, input LoginInput) (LoginOutput, error) {
	logger := s.Logger.WithField("service", "Auth.Login")

	user, err := s.Storage.Login(ctx, &storageAuth.LoginInput{
		Email: input.Email,
	})
	if err != nil {
		logger.Errorf("fail to login user: %v", err.Error())
		return LoginOutput{}, ErrFailedToLogin
	}

	if user == nil {
		logger.Errorf("fail to login user: %v", "account not found")
		return LoginOutput{}, ErrAccountNotFound
	}

	if !user.Verified {
		logger.Errorf("fail to login user: %v", "account hasn't been verified yet")
		return LoginOutput{}, ErrAccountNotVerified
	}

	if !util.ComparePasswords(user.Password, input.Password) {
		logger.Errorf("fail to login user: %v", "wrong password")
		return LoginOutput{}, ErrWrongPassword
	}

	return LoginOutput{
		UserId:   uuid.MustParse(user.UserID),
		FullName: user.FullName,
		Email:    user.Email,
	}, nil
}
