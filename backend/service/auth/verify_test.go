package auth_test

import (
	"context"
	"errors"
	"testing"

	"backend/service/auth"
	mockStorage "backend/storage/auth/mock"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestVerify(t *testing.T) {
	validUserID := uuid.New().String()

	tests := []struct {
		name      string
		input     auth.VerifyInput
		mockSetup func(*mockStorage.MockPostgresStore)
		wantOut   auth.VerifyOutput
		wantErr   error
	}{
		{
			name: "success",
			input: auth.VerifyInput{
				Email:             "user@example.com",
				VerificationToken: "ABC1234567",
			},
			mockSetup: func(store *mockStorage.MockPostgresStore) {
				store.EXPECT().
					CheckedUserForVerification(gomock.Any(), "user@example.com", "ABC1234567").
					Return(1, nil)

				store.EXPECT().
					VerifyUser(gomock.Any(), "user@example.com").
					Return(validUserID, nil)
			},
			wantOut: auth.VerifyOutput{UserID: uuid.MustParse(validUserID)},
		},
		{
			name: "CheckedUserForVerification storage error returns ErrFailedToVerify",
			input: auth.VerifyInput{
				Email:             "user@example.com",
				VerificationToken: "ABC1234567",
			},
			mockSetup: func(store *mockStorage.MockPostgresStore) {
				store.EXPECT().
					CheckedUserForVerification(gomock.Any(), "user@example.com", "ABC1234567").
					Return(0, errors.New("db error"))
			},
			wantErr: auth.ErrFailedToVerify,
		},
		{
			name: "user not found returns ErrAccountNotFound",
			input: auth.VerifyInput{
				Email:             "ghost@example.com",
				VerificationToken: "BADTOKEN00",
			},
			mockSetup: func(store *mockStorage.MockPostgresStore) {
				store.EXPECT().
					CheckedUserForVerification(gomock.Any(), "ghost@example.com", "BADTOKEN00").
					Return(0, nil)
			},
			wantErr: auth.ErrAccountNotFound,
		},
		{
			name: "VerifyUser storage error returns ErrFailedToVerify",
			input: auth.VerifyInput{
				Email:             "user@example.com",
				VerificationToken: "ABC1234567",
			},
			mockSetup: func(store *mockStorage.MockPostgresStore) {
				store.EXPECT().
					CheckedUserForVerification(gomock.Any(), "user@example.com", "ABC1234567").
					Return(1, nil)

				store.EXPECT().
					VerifyUser(gomock.Any(), "user@example.com").
					Return("", errors.New("update failed"))
			},
			wantErr: auth.ErrFailedToVerify,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockStorage.NewMockPostgresStore(ctrl)
			tc.mockSetup(store)

			svc := auth.New(logrus.New(), store)
			got, err := svc.Verify(context.Background(), tc.input)

			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.wantOut, got)
		})
	}
}
