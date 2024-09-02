package scotia

type ScotiaConfig struct {
	Mode             string `env_var:"SCOTIA_MODE,required"`
	BaseUrl          string `env_var:"SCOTIA_BASE_URL,required"`
	AuthUserName     string `env_var:"SCOTIA_BASIC_AUTH_USERNAME,required"`
	AuthPassword     string `env_var:"SCOTIA_BASIC_AUTH_PASSWORD,required"`
	ClientId         string `env_var:"SCOTIA_CLIENT_ID,required"`
	JWTKid           string `env_var:"SCOTIA_JWT_KID,required"`
	JWTAudience      string `env_var:"SCOTIA_JWT_AUDIENCE,required"`
	JWTExpiryMinutes int    `env_var:"SCOTIA_JWT_EXPIRY_MINUTES,required"`
	PrivateKeyDir    string `env_var:"SCOTIA_PRIVATE_KEY_DIR,required"`
	PublicKeyDir     string `env_var:"SCOTIA_Public_KEY_DIR,required"`
	Scope            string `env_var:"SCOTIA_SCOPE,required"`
	ProfileId        string `env_var:"SCOTIA_PROFILE_ID,required"`
	ApiKey           string `env_var:"SCOTIA_API_KEY,required"`
	CountryCode      string `env_var:"SCOTIA_COUNTRY_CODE,required"`
}
