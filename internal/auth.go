package internal

import "github.com/golang-jwt/jwt/v4"

type JwtClaims struct {
	ID int64 `json:"id"`

	jwt.RegisteredClaims
}
