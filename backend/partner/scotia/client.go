package scotia

import (
	cryptoRsa "crypto/rsa"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type rsa struct {
	privateKeyDir string
	signKey       *cryptoRsa.PrivateKey
}

func initRSA(config scotiaConfig) (*rsa, error) {
	//TODO: load rsa from config.
	return nil, nil
}

type ScotiaClientImpl struct {
	config scotiaConfig
	rsa    *rsa
}

func (s *ScotiaClientImpl) signJWT() (string, error) {
	claims := &jwt.RegisteredClaims{
		Subject:   s.config.ClientId,
		Audience:  []string{s.config.JWTAudience},
		Issuer:    s.config.ClientId,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(s.config.JWTExpiry) * time.Minute)),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	token.Header["kid"] = s.config.JWTKid
	ss, err := token.SignedString(s.rsa.signKey)
	if err != nil {
		return "", fmt.Errorf("ScotiaClientImpl JWT signature got error `%v`", err.Error())
	}
	return ss, nil
}
