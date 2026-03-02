package auth_storage

import (
	"context"
)

// sql queries
var (
	registerSQL = `
		INSERT INTO users (
			full_name,
			email,
			password,
			verification_token,
			phone_number
		) VALUES (
				$1, $2, $3, $4, $5
		) RETURNING user_id`
)

func (s *Storage) Register(ctx context.Context, reg Register) (string, error) {
	var id string
	if err := s.db.QueryRowContext(ctx, registerSQL,
		reg.FullName,
		reg.Email,
		reg.Password,
		reg.VerificationToken,
		reg.PhoneNumber,
	).Scan(&id); err != nil {
		return "", err
	}
	return id, nil
}
