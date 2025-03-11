package auth_storage

import (
	"context"
	"database/sql"
)

var (
	loginSQL = `
		SELECT 
			user_id,
			full_name,
			email,
			password,
			verified
		FROM users
		WHERE email = $1;
	`
)

func (s *Storage) Login(ctx context.Context, input *LoginInput) (*LoginOutput, error) {
	var output LoginOutput
	if err := s.db.QueryRowContext(ctx, loginSQL, input.Email).Scan(
		&output.UserID,
		&output.FullName,
		&output.Email,
		&output.Password,
		&output.Verified,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &output, nil
}
