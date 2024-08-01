package scotia

import (
	"fmt"
)

const (
	ConfigBaseUrl           = "SCOTIA_BASEURL"
	ConfigBasicAuthUsername = "SCOTIA_BASIC_AUTH_USERNAME"
	ConfigBasicAuthPassword = "SCOTIA_BASIC_AUTH_PASSWORD"
	ConfigClientId          = "SCOTIA_CLIENT_ID"
	ConfigJWTKid            = "SCOTIA_JWT_KID"
	ConfigJWTAudience       = "SCOTIA_JWT_AUDIENCE"
	ConfigJWTExpiry         = "SCOTIA_JWT_Expiry"
	ConfigPrivateKeyDir     = "SCOTIA_PRIVATE_KEY_DIR"
	ConfigPublicKeyDir      = "SCOTIA_Public_KEY_DIR"
	ConfigScope             = "SCOTIA_SCOPE"
)

func NewScotiaConfig() ScotiaConfig {
	m := make(map[string]interface{}, 15)
	return _scotiaConfig(m)
}

type ScotiaConfig interface {
	GetBaseUrl() string
	SetBaseUrl(u string)
	GetAuthUsername() string
	SetAuthUsername(u string)
	GetAuthPassword() string
	SetAuthPassword(u string)
	GetClientId() string
	SetClientId(u string)
	GetJWTKid() string
	SetJWTKid(u string)
	GetJWTAudience() string
	SetJWTAudience(u string)
	GetJWTExpiry() int
	SetJWTExpiry(u int)
	GetPrivateKeyDir() string
	SetPrivateKeyDir(u string)
	GetPublicKeyDir() string
	SetPublicKeyDir(u string)
	GetScope() string
	SetScope(u string)
}

type _scotiaConfig map[string]interface{}

func (c _scotiaConfig) GetBaseUrl() string {
	if val, ok := c[ConfigBaseUrl]; ok {
		v, k := val.(string)
		if k {
			return v
		}
		return ""
	}
	return ""
}

func (c _scotiaConfig) SetBaseUrl(u string) {

}

func (c _scotiaConfig) GetAuthUsername() string {
	if val, ok := c[ConfigBasicAuthUsername]; ok {
		v, k := val.(string)
		if k {
			return v
		}
		return ""
	}
	return ""
}

func (c _scotiaConfig) SetAuthUsername(u string) {

}

func (c _scotiaConfig) GetAuthPassword() string {
	if val, ok := c[ConfigBasicAuthPassword]; ok {
		v, k := val.(string)
		if k {
			return v
		}
		return ""
	}
	return ""
}

func (c _scotiaConfig) SetAuthPassword(u string) {

}

func (c _scotiaConfig) GetClientId() string {
	if val, ok := c[ConfigClientId]; ok {
		v, k := val.(string)
		if k {
			return v
		}
		return ""
	}
	return ""
}

func (c _scotiaConfig) SetClientId(u string) {

}

func (c _scotiaConfig) GetJWTKid() string {
	if val, ok := c[ConfigJWTKid]; ok {
		v, k := val.(string)
		if k {
			return v
		}
		return ""
	}
	return ""
}

func (c _scotiaConfig) SetJWTKid(u string) {

}

func (c _scotiaConfig) GetJWTAudience() string {
	if val, ok := c[ConfigJWTAudience]; ok {
		v, k := val.(string)
		if k {
			return v
		}
		return ""
	}
	return ""
}

func (c _scotiaConfig) SetJWTAudience(u string) {

}

func (c _scotiaConfig) GetJWTExpiry() int {
	if val, ok := c[ConfigJWTExpiry]; ok {
		v, k := val.(int)
		if k {
			return v
		}
		return 0
	}
	return 0
}

func (c _scotiaConfig) SetJWTExpiry(u int) {

}

func (c _scotiaConfig) GetScope() string {
	if val, ok := c[ConfigScope]; ok {
		v, k := val.(string)
		if k {
			return v
		}
		return ""
	}
	return ""
}

func (c _scotiaConfig) SetScope(u string) {

}

func (c _scotiaConfig) GetPrivateKeyDir() string {
	if val, ok := c[ConfigPrivateKeyDir]; ok {
		v, k := val.(string)
		if k {
			return v
		}
		return ""
	}
	return ""
}

func (c _scotiaConfig) SetPrivateKeyDir(u string) {

}

func (c _scotiaConfig) GetPublicKeyDir() string {
	if val, ok := c[ConfigPublicKeyDir]; ok {
		v, k := val.(string)
		if k {
			return v
		}
		return ""
	}
	return ""
}

func (c _scotiaConfig) SetPublicKeyDir(u string) {

}

func (c _scotiaConfig) String() string {
	return fmt.Sprintf("%v", c)
}
