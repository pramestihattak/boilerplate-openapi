//go:build integration

package auth_storage

import (
	"context"
	"log"
	"os"
	"testing"

	"backend/util"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

var testStorage *Storage

func TestMain(m *testing.M) {
	dbConfig := util.DBConfig{
		User:     getenv("TEST_DB_USER", "postgres"),
		Host:     getenv("TEST_DB_HOST", "localhost"),
		Port:     getenv("TEST_DB_PORT", "5438"),
		DBName:   getenv("TEST_DB_NAME", "postgres"),
		Password: getenv("TEST_DB_PASSWORD", "secret"),
		SSLMode:  "disable",
	}

	connStr, _ := util.NewDBStringFromDBConfig(dbConfig)

	logger := logrus.New()
	db, err := NewDbConn(logger, connStr)
	if err != nil {
		log.Printf("skipping integration tests: cannot connect to postgres: %v", err)
		os.Exit(0)
	}

	testStorage = &Storage{logger: logger, db: db}

	os.Exit(m.Run())
}

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

// seedUser inserts a user directly with full control over the verified state.
// Returns the inserted user_id.
func seedUser(t *testing.T, reg Register, verified bool) string {
	t.Helper()
	var id string
	err := testStorage.db.QueryRowContext(
		context.Background(),
		`INSERT INTO users (full_name, email, password, verification_token, phone_number, verified)
		 VALUES ($1, $2, $3, $4, $5, $6)
		 RETURNING user_id`,
		reg.FullName, reg.Email, reg.Password, reg.VerificationToken, reg.PhoneNumber, verified,
	).Scan(&id)
	require.NoError(t, err, "seedUser: failed to insert test user %q", reg.Email)
	return id
}

// cleanupByEmail deletes users by email. Call via defer for guaranteed cleanup.
func cleanupByEmail(t *testing.T, emails ...string) {
	t.Helper()
	for _, email := range emails {
		if _, err := testStorage.db.ExecContext(
			context.Background(),
			"DELETE FROM users WHERE email = $1",
			email,
		); err != nil {
			t.Logf("cleanup warning: failed to delete user %q: %v", email, err)
		}
	}
}
