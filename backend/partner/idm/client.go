package idm

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type IDMClient interface {
	Transfer(req IDMRequest) (*IDMResponse, error)
	GetConfigs() map[string]string
	SetConfig(key string, value string)
}

func NewIDMClient(config map[string]string) IDMClient {
	idmConfig := NewIDMConfigWithDefaultConfig(config)
	return &IDMClientImpl{
		config:     idmConfig,
		httpClient: &http.Client{},
	}
}

type IDMClientImpl struct {
	config     IDMConfig
	httpClient *http.Client
}

func (s *IDMClientImpl) GetConfigs() map[string]string {
	return map[string]string{}
}

func (s *IDMClientImpl) SetConfig(key string, value string) {
	s.config.SetConfig(key, value)
}

func (c *IDMClientImpl) Transfer(req IDMRequest) (*IDMResponse, error) {
	url := fmt.Sprintf("%s/account/transfer", c.config.GetBaseUrl())
	basicAuth := fmt.Sprintf("%s:%s", c.config.GetAuthUsername(), c.config.GetAuthPassword())
	basicAuth = base64.StdEncoding.EncodeToString([]byte(basicAuth))
	basicAuth = fmt.Sprintf("Basic %v", basicAuth)

	rawReqeust, err := json.Marshal(req)
	if err != nil {
		//Unlikely; Fatal
		return nil, err
	}

	r, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(rawReqeust))
	if err != nil {
		//Unlikely; Fatal
		return nil, err
	}
	r.Header.Add("Content-type", "application/json")
	r.Header.Add("Accept-Encoding", "UTF-8")
	r.Header.Add("Authorization", basicAuth)

	resp, err := c.httpClient.Do(r)
	if err != nil {
		//Unlikely; Fatal
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		//Unlikely; Fatal
		return nil, err
	}

	ret := &IDMResponse{
		ResponseCommon: ResponseCommon{
			StatusCode:  resp.StatusCode,
			RawRequest:  string(rawReqeust),
			RawResponse: string(body),
		},
	}

	if resp.StatusCode/100 != 2 {
		return ret, nil
	}

	derr := json.NewDecoder(bytes.NewBuffer(body)).Decode(ret)
	if derr != nil {
		//Unlikely; Fatal
		return nil, err
	}

	return ret, nil
}
