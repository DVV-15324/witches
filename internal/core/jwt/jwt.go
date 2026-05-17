package jwt

import (
	"context"
	"fmt"

	"time"

	"github.com/golang-jwt/jwt/v5"
	"sigs.k8s.io/kind/pkg/errors"
)

type Token struct {
	Token    string `json:"token"`
	ExpireAt int    `json:"expire_at"`
}
type TokenResponse struct {
	AccessToken Token  `json:"access_token"`
	RefeshToken *Token `json:"refesh_token"`
}

type JwtClaims struct {
	jwt.RegisteredClaims
}

func NewJwtServer(secretKey string, expireAt int) *JwtService {
	return &JwtService{
		secretKey: secretKey,
		expireAt:  expireAt,
	}
}

type JwtService struct {
	secretKey string
	expireAt  int
}

func (j *JwtService) IssueToken(cxt context.Context, sub string, tid string) *TokenResponse {
	now := time.Now()

	c := JwtClaims{
		jwt.RegisteredClaims{
			ID:        tid,
			Subject:   sub,
			NotBefore: jwt.NewNumericDate(now),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Duration(j.expireAt) * time.Second)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	signedToken, _ := token.SignedString([]byte(j.secretKey))

	return &TokenResponse{
		AccessToken: Token{
			Token:    signedToken,
			ExpireAt: j.expireAt,
		},
		RefeshToken: nil,
	}
}

func (j JwtService) ParseToken(ctx context.Context, tokenStr string) (*JwtClaims, error) {
	var claims JwtClaims
	token, err := jwt.ParseWithClaims(tokenStr, &claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.WithStack(fmt.Errorf("unexpected signing method: typ=%v, alg=%v", t.Header["typ"], t.Header["alg"]))
		}
		return []byte(j.secretKey), nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, errors.WithStack(fmt.Errorf("invalid token"))
	}
	return &claims, nil
}
