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
	var id string
	if err := s.db.QueryRowContext(ctx, verifyUserSQL, email).Scan(&id); err != nil {
		if err == sql.ErrNoRows {
			return "", sql.ErrNoRows
		}
		return "", err
	}
	return id, nil
}
