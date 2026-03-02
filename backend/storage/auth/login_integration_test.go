//go:build integration

package auth_storage

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLogin_Integration(t *testing.T) {
	tests := []struct {
		name        string
		seed        *Register
		seedVerified bool
		input       *LoginInput
		wantNil     bool
		verify      func(t *testing.T, got *LoginOutput)
	}{
		{
			name: "returns full user row for existing email",
			seed: &Register{
				FullName:          "Test User",
				Email:             "integ_login_found@test.com",
				Password:          "hashedpassword",
				VerificationToken: "TOKEN12345",
				PhoneNumber:       "08123456789",
			},
			seedVerified: true,
			input:        &LoginInput{Email: "integ_login_found@test.com"},
			verify: func(t *testing.T, got *LoginOutput) {
				require.NotNil(t, got)
				assert.Equal(t, "integ_login_found@test.com", got.Email)
				assert.Equal(t, "Test User", got.FullName)
				assert.Equal(t, "hashedpassword", got.Password)
				assert.True(t, got.Verified)
				assert.NotEmpty(t, got.UserID)
			},
		},
		{
			name: "returns verified=false for unverified user",
			seed: &Register{
				FullName:          "Unverified User",
				Email:             "integ_login_unverified@test.com",
				Password:          "hashedpassword",
				VerificationToken: "TOKEN12345",
				PhoneNumber:       "08123456789",
			},
			seedVerified: false,
			input:        &LoginInput{Email: "integ_login_unverified@test.com"},
			verify: func(t *testing.T, got *LoginOutput) {
				require.NotNil(t, got)
				assert.False(t, got.Verified)
			},
		},
		{
			name:    "returns nil for non-existent email",
			input:   &LoginInput{Email: "integ_login_notfound@test.com"},
			wantNil: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.seed != nil {
				seedUser(t, *tc.seed, tc.seedVerified)
				defer cleanupByEmail(t, tc.seed.Email)
			}

			got, err := testStorage.Login(context.Background(), tc.input)

			require.NoError(t, err)

			if tc.wantNil {
				assert.Nil(t, got)
				return
			}

			tc.verify(t, got)
		})
	}
}
