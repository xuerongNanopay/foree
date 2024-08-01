package scotia

import "crypto/rsa"

type ScotiaConfig struct {
	BaseUrl     string
	ClientId    string
	JWTKid      string
	JWTAudience string
	JWTExpiry   int // in minutes
	JWTSignKey  *rsa.PrivateKey
	Scope       string
}
