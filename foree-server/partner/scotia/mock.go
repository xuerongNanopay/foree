package scotia

func NewMockScotiaClient() ScotiaClient {
	return &ScotiaClientMock{
		config: ScotiaConfig{
			Mode: "mock",
		},
	}
}

type ScotiaClientMock struct {
	config ScotiaConfig
}

func (s *ScotiaClientMock) RequestPayment(req RequestPaymentRequest) (*RequestPaymentResponse, error) {
	return &RequestPaymentResponse{
		ResponseCommon: ResponseCommon{
			StatusCode: 200,
		},
		Data: RequestPaymentResponseData{
			PaymentId:               req.RequestData.EndToEndIdentification,
			ClearingSystemReference: "mock-clearing-system-reference",
		},
	}, nil
}
func (s *ScotiaClientMock) PaymentDetail(req PaymentDetailRequest) (*PaymentDetailResponse, error) {
	return &PaymentDetailResponse{
		ResponseCommon: ResponseCommon{
			StatusCode: 200,
		},
		PaymentDetail: PaymentDetailData{
			RequestForPaymentStatus: "ACCC",
		},
	}, nil
}
func (s *ScotiaClientMock) PaymentStatus(req PaymentStatusRequest) (*PaymentStatusResponse, error) {
	return &PaymentStatusResponse{
		ResponseCommon: ResponseCommon{
			StatusCode: 200,
		},
		PaymentStatuses: []PaymentStatusData{
			{
				GatewayUrl: "http://www.google.ca",
			},
		},
	}, nil
}
func (s *ScotiaClientMock) CancelPayment(req CancelPaymentRequest) (*CancelPaymentResponse, error) {
	return &CancelPaymentResponse{
		CancelStatus: CancelPaymentData{
			Status: "SUCCESS",
		},
	}, nil
}
func (s *ScotiaClientMock) GetConfigs() ScotiaConfig {
	return s.config
}
func (s *ScotiaClientMock) SetConfig(key string, value string) error {
	return nil
}
