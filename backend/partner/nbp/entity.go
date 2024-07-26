package nbp

type responseWrapper[R any] struct {
	RawRequest  string
	HttpStatus  int
	RawResponse int
	Data        R
}

type authenticate struct {
	ResponseCode    string `json:"ResponseCode"`
	ResponseMessage string `json:"ResponseMessage"`
	Token           string `json:"Token"`
	TokenExpiry     string `json:"Token_Expiry"`
}

type BankList struct {
	ResponseCode    string          `json:"ResponseCode"`
	ResponseMessage string          `json:"ResponseMessage"`
	Banklist        []BankListEntry `json:"banklist"`
}

type BankListEntry struct {
	BankName string `json:"bankName"`
	BankAbbr string `json:"bankAbbr"`
}

type HelloResponse responseWrapper[string]
type BankListResponse responseWrapper[BankList]

type authenticateRequest struct {
	AgencyCode int32
	UserName   string
	Password   string
}

// type authenticateRequest
