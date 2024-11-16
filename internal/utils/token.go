package utils

import (
	"context"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"

	"github.com/valek177/auth/internal/config"
	"github.com/valek177/auth/internal/model"
)

type Token interface {
	GenerateToken(_ context.Context, user *model.User) (string, error)
	VerifyToken(_ context.Context, token string) (*model.UserClaims, error)
}

type token struct {
	expTime time.Duration
	secret  []byte
}

func NewToken(cfg config.TokenConfig) *token {
	return &token{
		expTime: cfg.ExpTime(),
		secret:  cfg.Secret(),
	}
}

func (t *token) GenerateToken(_ context.Context, user *model.User) (string, error) {
	claims := model.UserClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(t.expTime).Unix(),
		},
		Username: user.Name,
		Role:     user.Role,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(t.secret)
}

func (t *token) VerifyToken(_ context.Context, token string) (*model.UserClaims, error) {
	fmt.Println("wea re in verify")
	tokenParsed, err := jwt.ParseWithClaims(
		token,
		&model.UserClaims{},
		func(token *jwt.Token) (interface{}, error) {
			_, ok := token.Method.(*jwt.SigningMethodHMAC)
			if !ok {
				return nil, errors.New("unexpected token signing method")
			}

			return t.secret, nil
		},
	)
	if err != nil {
		fmt.Println("err is ", err)
		return nil, errors.Errorf("invalid token: %s", err.Error())
	}
	fmt.Println("verify token", token, tokenParsed.Claims)

	claims, ok := tokenParsed.Claims.(*model.UserClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}
	fmt.Println("verify claims", claims)

	return claims, nil
}
