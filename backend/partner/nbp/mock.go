package nbp

func NewMockNBPClient() NBPClient {
	return &NBPClientMock{}
}

type NBPClientMock struct {
}

func (*NBPClientMock) Hello() (*HelloResponse, error) {
	return nil, nil
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
