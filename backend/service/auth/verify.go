package auth

import (
	"context"

	"github.com/google/uuid"
)

func (s *Auth) Verify(ctx context.Context, input VerifyInput) (VerifyOutput, error) {
	logger := s.Logger.WithField("service", "Auth.Verification")

	exist, err := s.Storage.CheckedUserForVerification(ctx, input.Email, input.VerificationToken)
	if err != nil {
		logger.Errorf("fail to verify user: %v", err.Error())
		return VerifyOutput{}, ErrFailedToVerify
	}

	if exist == 0 {
		logger.Errorf("fail to verify user: %v", "account not found")
		return VerifyOutput{}, ErrAccountNotFound
	}

	userID, err := s.Storage.VerifyUser(ctx, input.Email)
	if err != nil {
		logger.Errorf("fail to verify user: %v", err.Error())
		return VerifyOutput{}, ErrFailedToVerify
	}

	return VerifyOutput{
		UserID: uuid.MustParse(userID),
	}, nil
}
