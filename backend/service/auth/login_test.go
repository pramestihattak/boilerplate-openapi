package auth_test

import (
	"context"
	"errors"
	"testing"

	"backend/service/auth"
	storageAuth "backend/storage/auth"
	mockStorage "backend/storage/auth/mock"
	"backend/util"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestLogin(t *testing.T) {
	validUserID := uuid.New().String()
	correctHash, _ := util.HashAndSalt("password123")

	tests := []struct {
		name      string
		input     auth.LoginInput
		mockSetup func(*mockStorage.MockPostgresStore)
		wantOut   auth.LoginOutput
		wantErr   error
	}{
		{
			name: "success",
			input: auth.LoginInput{
				Email:    "user@example.com",
				Password: "password123",
			},
			mockSetup: func(store *mockStorage.MockPostgresStore) {
				store.EXPECT().
					Login(gomock.Any(), &storageAuth.LoginInput{Email: "user@example.com"}).
					Return(&storageAuth.LoginOutput{
						UserID:   validUserID,
						FullName: "Test User",
						Email:    "user@example.com",
						Password: correctHash,
						Verified: true,
					}, nil)
			},
			wantOut: auth.LoginOutput{
				UserID:   uuid.MustParse(validUserID),
				FullName: "Test User",
				Email:    "user@example.com",
			},
		},
		{
			name: "storage error returns ErrFailedToLogin",
			input: auth.LoginInput{
				Email:    "user@example.com",
				Password: "password123",
			},
			mockSetup: func(store *mockStorage.MockPostgresStore) {
				store.EXPECT().
					Login(gomock.Any(), &storageAuth.LoginInput{Email: "user@example.com"}).
					Return(nil, errors.New("db connection lost"))
			},
			wantErr: auth.ErrFailedToLogin,
		},
		{
			name: "nil user returns ErrAccountNotFound",
			input: auth.LoginInput{
				Email:    "ghost@example.com",
				Password: "password123",
			},
			mockSetup: func(store *mockStorage.MockPostgresStore) {
				store.EXPECT().
					Login(gomock.Any(), &storageAuth.LoginInput{Email: "ghost@example.com"}).
					Return(nil, nil)
			},
			wantErr: auth.ErrAccountNotFound,
		},
		{
			name: "unverified account returns ErrAccountNotVerified",
			input: auth.LoginInput{
				Email:    "user@example.com",
				Password: "password123",
			},
			mockSetup: func(store *mockStorage.MockPostgresStore) {
				store.EXPECT().
					Login(gomock.Any(), &storageAuth.LoginInput{Email: "user@example.com"}).
					Return(&storageAuth.LoginOutput{
						UserID:   validUserID,
						Email:    "user@example.com",
						Password: correctHash,
						Verified: false,
					}, nil)
			},
			wantErr: auth.ErrAccountNotVerified,
		},
		{
			name: "wrong password returns ErrWrongPassword",
			input: auth.LoginInput{
				Email:    "user@example.com",
				Password: "wrongpassword",
			},
			mockSetup: func(store *mockStorage.MockPostgresStore) {
				store.EXPECT().
					Login(gomock.Any(), &storageAuth.LoginInput{Email: "user@example.com"}).
					Return(&storageAuth.LoginOutput{
						UserID:   validUserID,
						Email:    "user@example.com",
						Password: correctHash,
						Verified: true,
					}, nil)
			},
			wantErr: auth.ErrWrongPassword,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockStorage.NewMockPostgresStore(ctrl)
			tc.mockSetup(store)

			svc := auth.New(logrus.New(), store)
			got, err := svc.Login(context.Background(), tc.input)

			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.wantOut, got)
		})
	}
}
