package auth_test

import (
	"context"
	"errors"
	"testing"

	"backend/service/auth"
	storageAuth "backend/storage/auth"
	mockStorage "backend/storage/auth/mock"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestRegister(t *testing.T) {
	validUserID := uuid.New().String()

	tests := []struct {
		name      string
		input     auth.RegisterInput
		mockSetup func(*mockStorage.MockPostgresStore)
		wantOut   auth.RegisterOutput
		wantErr   error
	}{
		{
			name: "success",
			input: auth.RegisterInput{
				FullName:    "Test User",
				Email:       "user@example.com",
				Password:    "password123",
				PhoneNumber: "08123456789",
			},
			mockSetup: func(store *mockStorage.MockPostgresStore) {
				store.EXPECT().
					UserExist(gomock.Any(), "user@example.com").
					Return(0, nil)

				store.EXPECT().
					Register(gomock.Any(), gomock.Any()).
					Do(func(_ context.Context, reg storageAuth.Register) {
						assert.Equal(t, "Test User", reg.FullName)
						assert.Equal(t, "user@example.com", reg.Email)
						assert.Equal(t, "08123456789", reg.PhoneNumber)
						assert.NotEmpty(t, reg.Password)
						assert.NotEqual(t, "password123", reg.Password, "password must be hashed")
						assert.Len(t, reg.VerificationToken, 10)
					}).
					Return(validUserID, nil)
			},
			wantOut: auth.RegisterOutput{UserID: uuid.MustParse(validUserID)},
		},
		{
			name: "UserExist storage error returns ErrFailedToRegister",
			input: auth.RegisterInput{
				FullName:    "Test User",
				Email:       "user@example.com",
				Password:    "password123",
				PhoneNumber: "08123456789",
			},
			mockSetup: func(store *mockStorage.MockPostgresStore) {
				store.EXPECT().
					UserExist(gomock.Any(), "user@example.com").
					Return(0, errors.New("db error"))
			},
			wantErr: auth.ErrFailedToRegister,
		},
		{
			name: "account already exists returns ErrAccountAlreadyExist",
			input: auth.RegisterInput{
				FullName:    "Test User",
				Email:       "existing@example.com",
				Password:    "password123",
				PhoneNumber: "08123456789",
			},
			mockSetup: func(store *mockStorage.MockPostgresStore) {
				store.EXPECT().
					UserExist(gomock.Any(), "existing@example.com").
					Return(1, nil)
			},
			wantErr: auth.ErrAccountAlreadyExist,
		},
		{
			name: "Register storage error returns ErrFailedToRegister",
			input: auth.RegisterInput{
				FullName:    "Test User",
				Email:       "user@example.com",
				Password:    "password123",
				PhoneNumber: "08123456789",
			},
			mockSetup: func(store *mockStorage.MockPostgresStore) {
				store.EXPECT().
					UserExist(gomock.Any(), "user@example.com").
					Return(0, nil)

				store.EXPECT().
					Register(gomock.Any(), gomock.Any()).
					Do(func(_ context.Context, reg storageAuth.Register) {
						assert.Equal(t, "Test User", reg.FullName)
						assert.Equal(t, "user@example.com", reg.Email)
						assert.Equal(t, "08123456789", reg.PhoneNumber)
					}).
					Return("", errors.New("insert failed"))
			},
			wantErr: auth.ErrFailedToRegister,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockStorage.NewMockPostgresStore(ctrl)
			tc.mockSetup(store)

			svc := auth.New(logrus.New(), store)
			got, err := svc.Register(context.Background(), tc.input)

			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.wantOut, got)
		})
	}
}
