package utils

import (
	"context"
	"fmt"

	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/pkg/errors"
)

type Token struct {
	Token    string `json:"token"`
	ExpireAt int    `json:"expire_at"`
}
type TokenResponse struct {
	AccessToken Token  `json:"access_token"`
	RefeshToken *Token `json:"refesh_token"`
}

type Role string

const (
	RoleAdmin Role = "admin"
	RoleUser  Role = "user"
)

type JwtClaims struct {
	Role string `json:"role"`
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

func (j *JwtService) IssueToken(ctx context.Context, sub string, tid string, role Role) (*TokenResponse, error) {

	now := time.Now()

	claims := JwtClaims{
		Role: string(role),
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        tid,
			Subject:   sub,
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(
				now.Add(time.Duration(j.expireAt) * time.Second),
			),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(j.secretKey))
	if err != nil {
		return nil, err
	}

	return &TokenResponse{
		AccessToken: Token{
			Token:    signedToken,
			ExpireAt: j.expireAt,
		},
	}, nil
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
		return nil, errors.WithStack(fmt.Errorf("Phiên của bạn đã hết vui lòng đăng nhập lại"))
	}
	if !token.Valid {
		return nil, errors.WithStack(fmt.Errorf("invalid token"))
	}
	return &claims, nil
}
