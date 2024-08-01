package scotia

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type ScotiaClientImpl struct {
	Config ScotiaConfig
}

func (s *ScotiaClientImpl) signJWT() (string, error) {
	claims := &jwt.RegisteredClaims{
		Subject:   s.Config.ClientId,
		Audience:  []string{s.Config.JWTAudience},
		Issuer:    s.Config.ClientId,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(s.Config.JWTExpiry) * time.Minute)),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	token.Header["kid"] = s.Config.JWTKid
	ss, err := token.SignedString(s.Config.JWTSignKey)
	if err != nil {
		return "", fmt.Errorf("ScotiaClientImpl JWT signature got error `%v`", err.Error())
	}
	return ss, nil
}
