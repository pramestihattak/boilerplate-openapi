package auth_storage

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

type Storage struct {
	logger logrus.FieldLogger
	db     *sqlx.DB

	PostgresStore
}

type PostgresStore interface {
	PostgresWriter
	PostgresReader
}

type PostgresWriter interface {
	Register(ctx context.Context, reg Register) (string, error)
	VerifyUser(ctx context.Context, email string) (string, error)
}

type PostgresReader interface {
	CheckedUserForVerification(ctx context.Context, email, verificationToken string) (int, error)
	UserExist(ctx context.Context, email string) (int, error)
	Login(ctx context.Context, input *LoginInput) (*LoginOutput, error)
}
