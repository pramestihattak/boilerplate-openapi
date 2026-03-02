package auth_storage

import (
	"context"
	"database/sql"
)

var getUserByIDSQL = `
	SELECT user_id, email, full_name, verified, phone_number
	FROM users
	WHERE user_id = $1
`

func (s *Storage) GetUserByID(ctx context.Context, id string) (*UserOutput, error) {
	var user UserOutput
	err := s.db.QueryRowContext(ctx, getUserByIDSQL, id).Scan(
		&user.UserID, &user.Email, &user.FullName, &user.Verified, &user.PhoneNumber,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}
