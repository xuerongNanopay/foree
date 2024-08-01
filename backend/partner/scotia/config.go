package scotia

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	ConfigBaseUrl           = "SCOTIA_BASE_URL"
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
	ShowConfigs() map[string]string
}

type _scotiaConfig map[string]interface{}

func NewScotiaConfig() ScotiaConfig {
	m := make(map[string]interface{}, 16)
	return _scotiaConfig(m)
}

func NewScotiaConfigWithDefaultConfig(configs map[string]string) ScotiaConfig {
	m := _scotiaConfig(make(map[string]interface{}, 16))
	if val, ok := configs[ConfigBaseUrl]; ok {
		m.SetBaseUrl(val)
	}
	if val, ok := configs[ConfigBasicAuthUsername]; ok {
		m.SetAuthUsername(val)
	}
	if val, ok := configs[ConfigBasicAuthPassword]; ok {
		m.SetAuthPassword(val)
	}
	if val, ok := configs[ConfigClientId]; ok {
		m.SetClientId(val)
	}
	if val, ok := configs[ConfigJWTKid]; ok {
		m.SetJWTKid(val)
	}
	if val, ok := configs[ConfigJWTAudience]; ok {
		m.SetJWTAudience(val)
	}
	if val, ok := configs[ConfigJWTExpiry]; ok {
		n, err := strconv.Atoi(val)
		if err != nil {
			//log?
			panic(err)
		}
		m.SetJWTExpiry(n)
	}
	if val, ok := configs[ConfigPrivateKeyDir]; ok {
		m.SetPrivateKeyDir(val)
	}
	if val, ok := configs[ConfigPublicKeyDir]; ok {
		m.SetPublicKeyDir(val)
	}
	if val, ok := configs[ConfigScope]; ok {
		m.SetScope(val)
	}
	return m
}

func (c _scotiaConfig) GetBaseUrl() string {
	return getStringConfig(c, ConfigBaseUrl)
}

func (c _scotiaConfig) SetBaseUrl(u string) {
	c[ConfigBaseUrl] = u
}

func (c _scotiaConfig) GetAuthUsername() string {
	return getStringConfig(c, ConfigBasicAuthUsername)
}

func (c _scotiaConfig) SetAuthUsername(u string) {
	c[ConfigBasicAuthUsername] = u
}

func (c _scotiaConfig) GetAuthPassword() string {
	return getStringConfig(c, ConfigBasicAuthPassword)
}

func (c _scotiaConfig) SetAuthPassword(u string) {
	c[ConfigBasicAuthPassword] = u
}

func (c _scotiaConfig) GetClientId() string {
	return getStringConfig(c, ConfigClientId)
}

func (c _scotiaConfig) SetClientId(u string) {
	c[ConfigClientId] = u
}

func (c _scotiaConfig) GetJWTKid() string {
	return getStringConfig(c, ConfigJWTKid)
}

func (c _scotiaConfig) SetJWTKid(u string) {
	c[ConfigJWTKid] = u
}

func (c _scotiaConfig) GetJWTAudience() string {
	return getStringConfig(c, ConfigJWTAudience)
}

func (c _scotiaConfig) SetJWTAudience(u string) {
	c[ConfigJWTAudience] = u
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
	c[ConfigJWTExpiry] = u
}

func (c _scotiaConfig) GetScope() string {
	return getStringConfig(c, ConfigScope)
}

func (c _scotiaConfig) SetScope(u string) {
	c[ConfigScope] = u
}

func (c _scotiaConfig) GetPrivateKeyDir() string {
	return getStringConfig(c, ConfigPrivateKeyDir)
}

func (c _scotiaConfig) SetPrivateKeyDir(u string) {
	c[ConfigPrivateKeyDir] = u
}

func (c _scotiaConfig) GetPublicKeyDir() string {
	return getStringConfig(c, ConfigPublicKeyDir)
}

func (c _scotiaConfig) SetPublicKeyDir(u string) {
	c[ConfigPublicKeyDir] = u
}

func getStringConfig(config _scotiaConfig, key string) string {
	if val, ok := config[key]; ok {
		return fmt.Sprintf("%v", val)
	}
	return ""
}

func (c _scotiaConfig) String() string {
	ret := []string{}
	for key, value := range c {
		ret = append(ret, fmt.Sprintf("%v:%v", key, value))
	}
	return strings.Join(ret, "\n")
}

func (c _scotiaConfig) ShowConfigs() map[string]string {
	ret := make(map[string]string, len(c))
	for key, value := range c {
		ret[key] = fmt.Sprintf("%v", value)
	}
	return ret
}
