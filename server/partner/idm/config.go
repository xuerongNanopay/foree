package idm

type IDMConfig struct {
	Mode         string `env_var:"IDM_MODE,required"`
	BaseUrl      string `env_var:"IDM_BASE_URL,required"`
	AuthUserName string `env_var:"IDM_BASIC_AUTH_USERNAME,required"`
	AuthPassword string `env_var:"IDM_BASIC_AUTH_PASSWORD,required"`
	HashingSalt  string `env_var:"IDM_HASHING_SALT,required"`
	Profile      string `env_var:"IDM_PROFILE,required"`
}
