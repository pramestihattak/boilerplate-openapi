package auth

import (
	storageAuth "backend/storage/auth"

	"github.com/sirupsen/logrus"
)

type Auth struct {
	Logger  *logrus.Logger
	Storage storageAuth.PostgresStore
}

func New(logger *logrus.Logger, storage storageAuth.PostgresStore) *Auth {
	return &Auth{
		Logger:  logger,
		Storage: storage,
	}
}
