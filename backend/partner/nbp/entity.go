package nbp

type responseWrapper[R any] struct {
	rawRequest  string
	httpStatus  int
	rawResponse int
	data        R
}

type BankList struct {
	responseCode    string          `json:"ResponseCode"`
	responseMessage string          `json:"ResponseMessage"`
	banklist        []BankListEntry `json:"banklist"`
}

type BankListEntry struct {
	bankName string `json:"bankName"`
	bankAbbr string `json:"bankAbbr"`
}

type HelloResponse responseWrapper[string]
type BankListResponse responseWrapper[BankList]
