package auth_storage

import (
	"context"
)

var (
	userExistSQL = `
		SELECT 
			COUNT(user_id)
		FROM users
		WHERE email = $1 LIMIT 1;
	`

	checkedUserVerificationExistSQL = `
		SELECT 
			COUNT(user_id)
		FROM users
		WHERE email = $1 and verification_token = $2 LIMIT 1;
	`
)

func (s *Storage) UserExist(ctx context.Context, email string) (int, error) {
	var exist int
	err := s.db.GetContext(ctx, &exist, userExistSQL, email)
	if err != nil {
		return 0, err
	}

	return exist, nil
}

func (s *Storage) CheckedUserForVerification(ctx context.Context, email, verificationToken string) (int, error) {
	var exist int
	err := s.db.GetContext(ctx, &exist, checkedUserVerificationExistSQL, email, verificationToken)
	if err != nil {
		return 0, err
	}

	return exist, nil
}
