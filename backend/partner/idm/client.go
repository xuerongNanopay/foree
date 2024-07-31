package idm

import "net/http"

type IDMClient interface {
	Transfer(req IDMRequest) (*IDMResponse, error)
}

func NewIDMClient(config IDMConfig) IDMClient {
	return &IDMClientImpl{
		Config:     config,
		httpClient: &http.Client{},
	}
}

type IDMClientImpl struct {
	Config     IDMConfig
	httpClient *http.Client
}

func (c *IDMClientImpl) Transfer(req IDMRequest) (*IDMResponse, error) {
	return nil, nil
}
