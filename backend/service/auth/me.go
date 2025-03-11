package auth

import "context"

func (a *Auth) Me(ctx context.Context, input MeInput) (MeOutput, error) {
	return MeOutput{}, nil
}
