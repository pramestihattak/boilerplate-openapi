//go:build integration

package auth_storage

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRegister_Integration(t *testing.T) {
	// Note: the users table has a composite unique(user_id, email) constraint,
	// not a standalone unique(email) constraint. Duplicate-email enforcement is
	// the responsibility of the service layer (via UserExist). Storage-level
	// tests only cover what the SQL layer actually enforces.
	tests := []struct {
		name string
		reg  Register
	}{
		{
			name: "inserts new user and returns a valid UUID",
			reg: Register{
				FullName:          "Test User",
				Email:             "integ_register_new@test.com",
				Password:          "hashedpassword",
				VerificationToken: "TOKEN12345",
				PhoneNumber:       "08123456789",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			defer cleanupByEmail(t, tc.reg.Email)

			id, err := testStorage.Register(context.Background(), tc.reg)

			require.NoError(t, err)
			_, parseErr := uuid.Parse(id)
			assert.NoError(t, parseErr, "returned ID should be a valid UUID")
		})
	}
}
