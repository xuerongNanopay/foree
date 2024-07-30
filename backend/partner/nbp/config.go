package nbp

type NBPConfig struct {
	BaseUrl              string
	AgencyCode           string
	Username             string
	Password             string
	AuthAttempts         int
	TokenExpiryThreshold int64
}
