package scotia

import "sync"

func NewMockScotiaClient() ScotiaClient {
	return &ScotiaClientMock{
		config: ScotiaConfig{
			Mode: "mock",
		},
	}
}

type ScotiaClientMock struct {
	config                ScotiaConfig
	cancelledTransactions sync.Map
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
	_, ok := s.cancelledTransactions.LoadAndDelete(req.PaymentId)
	if ok {
		return &PaymentDetailResponse{
			ResponseCommon: ResponseCommon{
				StatusCode: 200,
			},
			PaymentDetail: PaymentDetailData{
				RequestForPaymentStatus: "CANCELLED",
			},
		}, nil
	}
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
	s.cancelledTransactions.Store(req.PaymentId, true)
	return &CancelPaymentResponse{
		ResponseCommon: ResponseCommon{
			StatusCode: 200,
		},
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
