package nbp

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	ConfigBaseUrl              = "NBP_BASE_URL"
	ConfigBasicAuthUsername    = "NBP_BASIC_AUTH_USERNAME"
	ConfigBasicAuthPassword    = "NBP_BASIC_AUTH_PASSWORD"
	ConfigAgencyCode           = "NBP_AGENCY_CODE"
	ConfigAuthAttempts         = "NBP_CONFIG_AUTH_ATTEMPTS"
	ConfigTokenExpiryThreshold = "NBP_TOKEN_EXPIRY_THRESHOD"
)

type NBPConfig interface {
	GetBaseUrl() string
	SetBaseUrl(u string)
	GetAuthUsername() string
	SetAuthUsername(u string)
	GetAuthPassword() string
	SetAuthPassword(u string)
	GetAgencyCode() string
	SetAgencyCode(u string)
	GetAuthAttempts() int
	SetAuthAttempts(u int)
	GetTokenExpiryThreshold() int64
	SetTokenExpiryThreshold(u int64)
	ShowConfigs() map[string]string
}

type _nbpConfig map[string]interface{}

func NewNBPConfig() NBPConfig {
	m := make(map[string]interface{}, 16)
	return _nbpConfig(m)
}

func NewNBPConfigWithDefaultConfig(configs map[string]string) NBPConfig {
	m := _nbpConfig(make(map[string]interface{}, len(configs)))
	if val, ok := configs[ConfigBaseUrl]; ok {
		m.SetBaseUrl(val)
	}
	if val, ok := configs[ConfigBasicAuthUsername]; ok {
		m.SetAuthUsername(val)
	}
	if val, ok := configs[ConfigBasicAuthPassword]; ok {
		m.SetAuthPassword(val)
	}
	if val, ok := configs[ConfigAgencyCode]; ok {
		m.SetAgencyCode(val)
	}
	if val, ok := configs[ConfigAuthAttempts]; ok {
		n, err := strconv.Atoi(val)
		if err != nil {
			//log?
			panic(err)
		}
		m.SetAuthAttempts(n)
	}
	if val, ok := configs[ConfigTokenExpiryThreshold]; ok {
		n, err := strconv.Atoi(val)
		if err != nil {
			//log?
			panic(err)
		}
		m.SetTokenExpiryThreshold(int64(n))
	}
	return m
}

func (c _nbpConfig) GetBaseUrl() string {
	return getStringConfig(c, ConfigBaseUrl)
}

func (c _nbpConfig) SetBaseUrl(u string) {
	c[ConfigBaseUrl] = u
}

func (c _nbpConfig) GetAuthUsername() string {
	return getStringConfig(c, ConfigBasicAuthUsername)
}

func (c _nbpConfig) SetAuthUsername(u string) {
	c[ConfigBasicAuthUsername] = u
}

func (c _nbpConfig) GetAuthPassword() string {
	return getStringConfig(c, ConfigBasicAuthPassword)
}

func (c _nbpConfig) SetAuthPassword(u string) {
	c[ConfigBasicAuthPassword] = u
}

func (c _nbpConfig) GetAgencyCode() string {
	return getStringConfig(c, ConfigAgencyCode)
}

func (c _nbpConfig) SetAgencyCode(u string) {
	c[ConfigAgencyCode] = u
}

func (c _nbpConfig) GetAuthAttempts() int {
	if val, ok := c[ConfigAuthAttempts]; ok {
		v, k := val.(int)
		if k {
			return v
		}
		return 0
	}
	return 0

}

func (c _nbpConfig) SetAuthAttempts(u int) {
	c[ConfigAuthAttempts] = u
}

func (c _nbpConfig) GetTokenExpiryThreshold() int64 {
	if val, ok := c[ConfigTokenExpiryThreshold]; ok {
		v, k := val.(int64)
		if k {
			return v
		}
		return 0
	}
	return 0

}

func (c _nbpConfig) SetTokenExpiryThreshold(u int64) {
	c[ConfigTokenExpiryThreshold] = u
}

func getStringConfig(config _nbpConfig, key string) string {
	if val, ok := config[key]; ok {
		return fmt.Sprintf("%v", val)
	}
	return ""
}

func (c _nbpConfig) String() string {
	ret := []string{}
	for key, value := range c {
		ret = append(ret, fmt.Sprintf("%v:%v", key, value))
	}
	return strings.Join(ret, "\n")
}

func (c _nbpConfig) ShowConfigs() map[string]string {
	ret := make(map[string]string, len(c))
	for key, value := range c {
		ret[key] = fmt.Sprintf("%v", value)
	}
	return ret
}
