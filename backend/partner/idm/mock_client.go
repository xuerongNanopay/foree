package idm

func NewMockIDMClient() IDMClient {

}

type IDMClientMock struct {
}

func (c *IDMClientMock) Transfer(req IDMRequest) (*IDMResponse, error) {
	return &IDMResponse{
		FraudEvaluationResult: ResultStatusAccept,
	}, nil
}
