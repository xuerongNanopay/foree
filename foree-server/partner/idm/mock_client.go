package idm

func NewMockIDMClient() IDMClient {
	return &IDMClientMock{
		config: IDMConfig{
			Mode: "mock",
		},
	}
}

type IDMClientMock struct {
	config IDMConfig
}

func (s *IDMClientMock) GetConfigs() IDMConfig {
	return s.config
}

func (s *IDMClientMock) SetConfig(key string, value string) error {
	return nil
}

func (c *IDMClientMock) Transfer(req IDMRequest) (*IDMResponse, error) {
	return &IDMResponse{
		ResponseCommon: ResponseCommon{
			StatusCode: 200,
		},
		FraudEvaluationResult: ResultStatusAccept,
	}, nil
}
