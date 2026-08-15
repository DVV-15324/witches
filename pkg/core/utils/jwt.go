package utils

import (
	"context"
	"fmt"

	"time"

	"github.com/DVV-15324/witches/cmd/utils"
	"github.com/golang-jwt/jwt/v5"
	"github.com/pkg/errors"
)

// Định nghĩa các lỗi cụ thể
var (
	ErrTokenExpired     = errors.New("token has expired") //Tín hiệu để FE biết để lấy token mới
	ErrInvalidSignature = errors.New("invalid token signature")
	ErrMalformedToken   = errors.New("malformed token")
	ErrInvalidAlgorithm = errors.New("invalid algorithm")
	ErrTokenNotYetValid = errors.New("token not yet valid")
	ErrInvalidToken     = errors.New("invalid token")
)

type Token struct {
	Token    string `json:"token"`
	ExpireAt int64  `json:"expire_at"`
}
type TokenResponse struct {
	AccessToken  Token `json:"access_token"`
	RefreshToken Token `json:"refresh_token"`
}

type JwtClaims struct {
	jwt.RegisteredClaims
}

type JwtService struct {
	config *utils.Config
}

func NewJwtService(config *utils.Config) *JwtService {
	return &JwtService{
		config: config,
	}
}

func (j *JwtService) IssueAccessToken(ctx context.Context, sub string, tid string) (*Token, error) {
	now := time.Now()

	claims := JwtClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        tid,
			Subject:   sub,
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(
				now.Add(time.Duration(j.config.AccessTokenTTL) * time.Second),
			),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(j.config.JWTSecret))
	if err != nil {
		return nil, err
	}

	return &Token{Token: signedToken, ExpireAt: j.config.AccessTokenTTL}, nil
}

func (j *JwtService) IssueRefreshToken(ctx context.Context, sub string, tid string) (*Token, error) {
	now := time.Now()

	claims := JwtClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        tid,
			Subject:   sub,
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(
				now.Add(time.Duration(j.config.RefreshTokenTTL) * time.Second),
			),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(j.config.JWTSecret))
	if err != nil {
		return nil, err
	}

	return &Token{Token: signedToken, ExpireAt: j.config.RefreshTokenTTL}, nil
}
func (j *JwtService) IssueTokenPair(ctx context.Context, sub string, tid string) (*TokenResponse, error) {
	accessToken, err := j.IssueAccessToken(ctx, sub, tid)
	if err != nil {
		return nil, err
	}

	refreshToken, err := j.IssueRefreshToken(ctx, sub, tid)
	if err != nil {
		return nil, err
	}

	return &TokenResponse{
		AccessToken:  *accessToken,
		RefreshToken: *refreshToken,
	}, nil
}
func (j *JwtService) ParseToken(ctx context.Context, tokenStr string) (*JwtClaims, error) {
	var claims JwtClaims

	token, err := jwt.ParseWithClaims(tokenStr, &claims, func(t *jwt.Token) (interface{}, error) {
		// Kiểm tra algorithm
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(j.config.JWTSecret), nil
	})

	if err != nil {
		switch {
		case errors.Is(err, jwt.ErrTokenExpired):
			return nil, ErrTokenExpired
		case errors.Is(err, jwt.ErrTokenSignatureInvalid):
			return nil, ErrInvalidSignature
		case errors.Is(err, jwt.ErrTokenMalformed):
			return nil, ErrMalformedToken
		case errors.Is(err, jwt.ErrTokenNotValidYet):
			return nil, ErrTokenNotYetValid
		case errors.Is(err, jwt.ErrTokenUnverifiable):
			return nil, ErrInvalidToken
		default:
			// Kiểm tra lỗi algorithm
			if errors.Is(err, jwt.ErrTokenInvalidClaims) {
				return nil, ErrInvalidToken
			}
			return nil, errors.Wrap(err, "parse token failed")
		}
	}

	// Kiểm tra token valid
	if !token.Valid {
		return nil, ErrInvalidToken
	}

	//  Kiểm tra thêm 1 lần nữa (phòng trường hợp)
	if claims.ExpiresAt != nil && claims.ExpiresAt.Before(time.Now()) {
		return nil, ErrTokenExpired
	}

	//  Kiểm tra NotBefore
	if claims.NotBefore != nil && claims.NotBefore.After(time.Now()) {
		return nil, ErrTokenNotYetValid
	}

	return &claims, nil
}
