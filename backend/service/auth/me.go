package auth

import "context"

func (a *Auth) Me(ctx context.Context, input MeInput) (MeOutput, error) {
	logger := a.Logger.WithField("service", "Auth.Me")

	user, err := a.Storage.GetUserByID(ctx, input.UserID)
	if err != nil {
		logger.Errorf("fail to get user: %v", err.Error())
		return MeOutput{}, ErrFailedToGetUser
	}
	if user == nil {
		return MeOutput{}, ErrAccountNotFound
	}

	return MeOutput{
		UserID:      user.UserID,
		Email:       user.Email,
		FullName:    user.FullName,
		Verified:    user.Verified,
		PhoneNumber: user.PhoneNumber,
	}, nil
}
