package token

import (
	"crypto/rsa"
	"fmt"
	"github.com/golang-jwt/jwt"
)

type JWTVerifier struct {
	PublicKey *rsa.PublicKey
}

func (v *JWTVerifier) Verify(token string) (string, error) {
	claims, err := jwt.ParseWithClaims(token, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return v.PublicKey, nil
	})
	if err != nil {
		return "", fmt.Errorf("cannot parse token:%v", err)
	}

	if !claims.Valid {
		return "", fmt.Errorf("token not valid")
	}

	clm, ok := claims.Claims.(*jwt.StandardClaims)
	if !ok {
		return "", fmt.Errorf("token claim is not standard claim")
	}
	if err := clm.Valid(); err != nil {
		return "", fmt.Errorf("claim not valid:%v", err)
	}
	return clm.Subject, nil
}
