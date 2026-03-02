//go:build integration

package auth_storage

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserExist_Integration(t *testing.T) {
	tests := []struct {
		name      string
		seed      *Register
		email     string
		wantCount int
	}{
		{
			name: "returns 1 when user exists",
			seed: &Register{
				FullName:          "Test User",
				Email:             "integ_userexist_found@test.com",
				Password:          "hashedpassword",
				VerificationToken: "TOKEN12345",
				PhoneNumber:       "08123456789",
			},
			email:     "integ_userexist_found@test.com",
			wantCount: 1,
		},
		{
			name:      "returns 0 when user does not exist",
			email:     "integ_userexist_notfound@test.com",
			wantCount: 0,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.seed != nil {
				seedUser(t, *tc.seed, false)
				defer cleanupByEmail(t, tc.seed.Email)
			}

			count, err := testStorage.UserExist(context.Background(), tc.email)

			require.NoError(t, err)
			assert.Equal(t, tc.wantCount, count)
		})
	}
}

func TestCheckedUserForVerification_Integration(t *testing.T) {
	seed := Register{
		FullName:          "Test User",
		Email:             "integ_checkverify@test.com",
		Password:          "hashedpassword",
		VerificationToken: "VALIDTOKEN",
		PhoneNumber:       "08123456789",
	}

	// Seed once — all subtests are read-only, so shared data is safe.
	seedUser(t, seed, false)
	defer cleanupByEmail(t, seed.Email)

	tests := []struct {
		name              string
		email             string
		verificationToken string
		wantCount         int
	}{
		{
			name:              "returns 1 for correct email and token",
			email:             "integ_checkverify@test.com",
			verificationToken: "VALIDTOKEN",
			wantCount:         1,
		},
		{
			name:              "returns 0 for wrong token",
			email:             "integ_checkverify@test.com",
			verificationToken: "WRONGTOKEN",
			wantCount:         0,
		},
		{
			name:              "returns 0 for non-existent email",
			email:             "integ_checkverify_ghost@test.com",
			verificationToken: "VALIDTOKEN",
			wantCount:         0,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			count, err := testStorage.CheckedUserForVerification(context.Background(), tc.email, tc.verificationToken)

			require.NoError(t, err)
			assert.Equal(t, tc.wantCount, count)
		})
	}
}
