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
	txn, err := s.db.Begin()
	if err != nil {
		return "", err
	}

	var id string
	if err := txn.QueryRowContext(ctx, registerSQL,
		reg.FullName,
		reg.Email,
		reg.Password,
		reg.VerificationToken,
		reg.PhoneNumber,
	).Scan(&id); err != nil {
		if err := txn.Rollback(); err != nil {
			return "", err
		}
		return "", err
	}

	if err := txn.Commit(); err != nil {
		return "", err
	}
	return id, nil
}
