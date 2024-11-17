package model

import "github.com/dgrijalva/jwt-go"

// UserClaims is a model for user claims
type UserClaims struct {
	jwt.StandardClaims
	Username string `json:"username"`
	Role     string `json:"role"`
}
