package idm

func NewMockIDMClient() IDMClient {
	return &IDMClientImpl{}
}

type IDMClientMock struct {
}

func (c *IDMClientMock) Transfer(req IDMRequest) (*IDMResponse, error) {
	return &IDMResponse{
		FraudEvaluationResult: ResultStatusAccept,
	}, nil
}
