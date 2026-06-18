package auth

import (
	c_jwt "arc-golang/internal/utils"
	"context"
)

func (a *UsecaseAuth) UsecaseIntrospectToken(ctx context.Context, accessToken string) (*c_jwt.JwtClaims, error) {
	claims, error := a.Jwt.ParseToken(ctx, accessToken)
	if error != nil {
		return nil, error
	}
	return claims, nil
}
