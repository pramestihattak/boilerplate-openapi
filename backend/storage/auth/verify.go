package auth_storage

import (
	"context"
	"database/sql"
)

var (
	verifyUserSQL = `
		UPDATE users
			SET verified = true, verification_token = ''
		WHERE email = $1
		RETURNING user_id
	`
)

func (s *Storage) VerifyUser(ctx context.Context, email string) (string, error) {
	txn, err := s.db.Begin()
	if err != nil {
		return "", err
	}

	var id string
	if err := txn.QueryRowContext(ctx, verifyUserSQL,
		email,
	).Scan(&id); err != nil {
		if err == sql.ErrNoRows {
			return "", nil
		}
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
