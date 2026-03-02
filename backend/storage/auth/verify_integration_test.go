//go:build integration

package auth_storage

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestVerifyUser_Integration(t *testing.T) {
	tests := []struct {
		name      string
		seed      *Register
		email     string
		wantEmpty bool
		verify    func(t *testing.T, got string)
	}{
		{
			name: "sets verified=true, clears token, and returns valid UUID",
			seed: &Register{
				FullName:          "Test User",
				Email:             "integ_verify_success@test.com",
				Password:          "hashedpassword",
				VerificationToken: "TOKEN12345",
				PhoneNumber:       "08123456789",
			},
			email: "integ_verify_success@test.com",
			verify: func(t *testing.T, got string) {
				_, parseErr := uuid.Parse(got)
				assert.NoError(t, parseErr, "returned ID should be a valid UUID")

				var verified bool
				var token string
				err := testStorage.db.QueryRowContext(
					context.Background(),
					"SELECT verified, verification_token FROM users WHERE email = $1",
					"integ_verify_success@test.com",
				).Scan(&verified, &token)
				require.NoError(t, err)
				assert.True(t, verified)
				assert.Empty(t, token)
			},
		},
		{
			name:      "non-existent email returns empty string without error",
			email:     "integ_verify_notfound@test.com",
			wantEmpty: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.seed != nil {
				seedUser(t, *tc.seed, false)
				defer cleanupByEmail(t, tc.seed.Email)
			}

			got, err := testStorage.VerifyUser(context.Background(), tc.email)

			require.NoError(t, err)

			if tc.wantEmpty {
				assert.Empty(t, got)
				return
			}

			tc.verify(t, got)
		})
	}
}
