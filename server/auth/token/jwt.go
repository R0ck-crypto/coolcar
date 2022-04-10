package token

import (
	"crypto/rsa"
	"github.com/golang-jwt/jwt"
	"time"
)

type JwtTokeGen struct {
	privateKey *rsa.PrivateKey
	issuer     string
	nowFunc    func() time.Time
}

func NewJwtTokeGen(issuer string, privateKey *rsa.PrivateKey) *JwtTokeGen {
	return &JwtTokeGen{
		privateKey: privateKey,
		issuer:     issuer,
		nowFunc:    time.Now,
	}
}

func (j *JwtTokeGen) GenerateToken(accountID string, expire time.Duration) (string, error) {
	nowSec := j.nowFunc().Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodRS512, jwt.StandardClaims{
		Issuer:    j.issuer,
		IssuedAt:  nowSec,
		ExpiresAt: nowSec + int64(expire.Seconds()),
		Subject:   accountID,
	})
	return token.SignedString(j.privateKey)
}
