package idm

import (
	"fmt"
	"strings"
)

const (
	ConfigBaseUrl           = "IDM_BASE_URL"
	ConfigBasicAuthUsername = "IDM_BASIC_AUTH_USERNAME"
	ConfigBasicAuthPassword = "IDM_BASIC_AUTH_PASSWORD"
	ConfigHashingSalt       = "IDM_HASHING_SALT"
	ConfigProfile           = "IDM_PROFILE"
)

type IDMConfig interface {
	GetBaseUrl() string
	SetBaseUrl(u string)
	GetAuthUsername() string
	SetAuthUsername(u string)
	GetAuthPassword() string
	SetAuthPassword(u string)
	GetHashingSalt() string
	SetHashingSalt(u string)
	GetProfile() string
	SetProfile(u string)
	ShowConfigs() map[string]string
}

type _idmConfig map[string]interface{}

func NewIDMConfig() IDMConfig {
	m := make(map[string]interface{}, 16)
	return _idmConfig(m)
}

func NewIDMConfigWithDefaultConfig(configs map[string]string) IDMConfig {
	m := _idmConfig(make(map[string]interface{}, len(configs)))
	if val, ok := configs[ConfigBaseUrl]; ok {
		m.SetBaseUrl(val)
	}
	if val, ok := configs[ConfigBasicAuthUsername]; ok {
		m.SetAuthUsername(val)
	}
	if val, ok := configs[ConfigBasicAuthPassword]; ok {
		m.SetAuthPassword(val)
	}
	if val, ok := configs[ConfigHashingSalt]; ok {
		m.SetHashingSalt(val)
	}
	if val, ok := configs[ConfigProfile]; ok {
		m.SetProfile(val)
	}
	return m
}

func (c _idmConfig) GetBaseUrl() string {
	return getStringConfig(c, ConfigBaseUrl)
}

func (c _idmConfig) SetBaseUrl(u string) {
	c[ConfigBaseUrl] = u
}

func (c _idmConfig) GetAuthUsername() string {
	return getStringConfig(c, ConfigBasicAuthUsername)
}

func (c _idmConfig) SetAuthUsername(u string) {
	c[ConfigBasicAuthUsername] = u
}

func (c _idmConfig) GetAuthPassword() string {
	return getStringConfig(c, ConfigBasicAuthPassword)
}

func (c _idmConfig) SetAuthPassword(u string) {
	c[ConfigBasicAuthPassword] = u
}

func (c _idmConfig) GetHashingSalt() string {
	return getStringConfig(c, ConfigHashingSalt)
}

func (c _idmConfig) SetHashingSalt(u string) {
	c[ConfigHashingSalt] = u
}

func (c _idmConfig) GetProfile() string {
	return getStringConfig(c, ConfigProfile)
}

func (c _idmConfig) SetProfile(u string) {
	c[ConfigProfile] = u
}

func getStringConfig(config _idmConfig, key string) string {
	if val, ok := config[key]; ok {
		return fmt.Sprintf("%v", val)
	}
	return ""
}

func (c _idmConfig) String() string {
	ret := []string{}
	for key, value := range c {
		ret = append(ret, fmt.Sprintf("%v:%v", key, value))
	}
	return strings.Join(ret, "\n")
}

func (c _idmConfig) ShowConfigs() map[string]string {
	ret := make(map[string]string, len(c))
	for key, value := range c {
		ret[key] = fmt.Sprintf("%v", value)
	}
	return ret
}
