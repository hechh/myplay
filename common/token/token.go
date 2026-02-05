package token

import (
	"myplay/common/pb"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/hechh/library/crypto"
)

type SessionToken struct {
	jwt.RegisteredClaims
	*pb.SessionData
}

func ParseToken(str string, secret string) (*pb.SessionData, error) {
	tok := &SessionToken{}
	if err := crypto.JwtDecrypto(str, secret, tok); err != nil {
		return nil, err
	}
	return tok.SessionData, nil
}

func GenToken(data *pb.SessionData, secret string) (string, error) {
	now := time.Now()
	item := &SessionToken{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   "myplay/auth",
			ExpiresAt: jwt.NewNumericDate(now.Add(30 * time.Second)),
			NotBefore: jwt.NewNumericDate(now),
			IssuedAt:  jwt.NewNumericDate(now),
		},
		SessionData: data,
	}
	return crypto.JwtEncrypto(item, secret)
}
