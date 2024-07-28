package nbp

import "net/http"

func NewMockNBPClient() NBPClient {
	return &NBPClientMock{}
}

type NBPClientMock struct {
}

func (*NBPClientMock) Hello() (*HelloResponse, error) {
	r := &HelloResponse{
		HttpStatus:  http.StatusOK,
		RawResponse: "Welcome, NBP E-Remittance API",
		Data:        "Welcome, NBP E-Remittance API",
	}
	return r, nil
}

func (*NBPClientMock) BankList() (*BankListResponse, error) {
	return nil, nil
}

func (*NBPClientMock) AccountEnquiry(r AccountEnquiryRequest) (*AccountEnquiryResponse, error) {
	return nil, nil
}

func (*NBPClientMock) LoadRemittance(r LoadRemittanceRequest) (*LoadRemittanceResponse, error) {
	return nil, nil
}

func (*NBPClientMock) TransactionStatusByIds(r TransactionStatusByIdsRequest) (*TransactionStatusByIdsResponse, error) {
	return nil, nil
}

func (*NBPClientMock) TransactionStatusByDate(r TransactionStatusByDateRequest) (*TransactionStatusByDateResponse, error) {
	return nil, nil
}
func (*NBPClientMock) CancelTransaction(r CancelTransactionRequest) (*CancelTransactionResponse, error) {
	return nil, nil
}
