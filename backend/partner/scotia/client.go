package scotia

import (
	cryptoRsa "crypto/rsa"
	"fmt"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type ScotiaClient interface {
	GetConfigs() map[string]string
	SetConfig(key string, value string)
}

type tokenData struct {
	token          string
	rawTokenExpiry string
	tokenExpiry    *time.Time
}

type rsa struct {
	privateKeyDir string
	signKey       *cryptoRsa.PrivateKey
}

func initRSA(config ScotiaConfig) (*rsa, error) {
	//TODO: load rsa from config.
	return nil, nil
}

func NewScotiaClientImpl(configs map[string]string) ScotiaClient {
	scotiaConfig := NewScotiaConfigWithDefaultConfig(configs)

	return &scotiaClientImpl{
		config: scotiaConfig,
	}
}

type scotiaClientImpl struct {
	config ScotiaConfig
	rsa    *rsa
	auth   *tokenData
	mu     sync.Mutex
}

func (s *scotiaClientImpl) GetConfigs() map[string]string {
	return s.config.ShowConfigs()
}

func (s *scotiaClientImpl) SetConfig(key string, value string) {
	s.config.SetConfig(key, value)
}

func (s *scotiaClientImpl) signJWT() (string, error) {
	claims := &jwt.RegisteredClaims{
		Subject:   s.config.GetClientId(),
		Audience:  []string{s.config.GetJWTAudience()},
		Issuer:    s.config.GetClientId(),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(s.config.GetJWTExpiryMinutes()) * time.Minute)),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	token.Header["kid"] = s.config.GetJWTKid()
	ss, err := token.SignedString(s.rsa.signKey)
	if err != nil {
		return "", fmt.Errorf("ScotiaClientImpl JWT signature got error `%v`", err.Error())
	}
	return ss, nil
}
