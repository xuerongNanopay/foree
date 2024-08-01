package idm

func NewMockIDMClient() IDMClient {
	return &IDMClientImpl{}
}

type IDMClientMock struct {
}

func (s *IDMClientMock) GetConfigs() map[string]string {
	return map[string]string{}
}

func (s *IDMClientMock) SetConfig(key string, value string) {
}

func (c *IDMClientMock) Transfer(req IDMRequest) (*IDMResponse, error) {
	return &IDMResponse{
		FraudEvaluationResult: ResultStatusAccept,
	}, nil
}
