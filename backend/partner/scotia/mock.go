package scotia

func NewMockScotiaClient() ScotiaClient {
	return &ScotiaClientMock{}
}

type ScotiaClientMock struct {
	config ScotiaConfig
}

func (s *ScotiaClientMock) RequestPayment(req RequestPaymentRequest) (*RequestPaymentResponse, error) {
	return nil, nil
}
func (s *ScotiaClientMock) PaymentDetail(req PaymentDetailRequest) (*PaymentDetailResponse, error) {
	return nil, nil
}
func (s *ScotiaClientMock) PaymentStatus(req PaymentStatusRequest) (*PaymentStatusResponse, error) {
	return nil, nil
}
func (s *ScotiaClientMock) CancelPayment(req CancelPaymentRequest) (*CancelPaymentResponse, error) {
	return nil, nil
}
func (s *ScotiaClientMock) GetConfigs() map[string]string {
	s.config.ShowConfigs()
}
func (s *ScotiaClientMock) SetConfig(key string, value string) {
	s.config.SetConfig(key, value)
}
