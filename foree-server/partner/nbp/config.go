package nbp

type NBPConfig struct {
	Mode                         string `env_var:"NBP_MODE,required"`
	BaseUrl                      string `env_var:"NBP_BASE_URL,required"`
	AuthUserName                 string `env_var:"NBP_BASIC_AUTH_USERNAME,required"`
	AuthPassword                 string `env_var:"NBP_BASIC_AUTH_PASSWORD,required"`
	AgencyCode                   string `env_var:"NBP_AGENCY_CODE,required"`
	TokenExpiryThresholdInSecond int64  `env_var:"NBP_TOKEN_EXPIRY_THRESHOD,default=60"`
}
