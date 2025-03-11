package auth

import "context"

type AuthService interface {
	Me(ctx context.Context, input MeInput) (MeOutput, error)
	Login(ctx context.Context, input LoginInput) (LoginOutput, error)
	Register(ctx context.Context, input RegisterInput) (RegisterOutput, error)
	Verify(ctx context.Context, input VerifyInput) (VerifyOutput, error)
}
