package nbp

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type authCache struct {
	token          string
	rawTokenExpiry string
	tokenExpiry    time.Time
}

type NBPClient interface {
	Hello() (*HelloResponse, error)
	BankList() (*BankListResponse, error)
	AccountEnquiry(AccountEnquiryRequest) (*AccountEnquiryResponse, error)
	LoadRemittance(LoadRemittanceRequest) (*LoadRemittanceResponse, error)
	TransactionStatusByIds(TransactionStatusByIdsRequest) (*TransactionStatusByIdsResponse, error)
	TransactionStatusByDate(TransactionStatusByDateRequest) (*TransactionStatusByDateResponse, error)
	CancelTransaction(CancelTransactionRequest) (*CancelTransactionResponse, error)
}

func NewNBPClient(config NBPConfig) NBPClient {
	return &NBPClientImpl{
		Config: config,
		httpClient: &http.Client{
			Timeout: 4 * time.Minute, // At least 3 minutes
		},
	}
}

type NBPClientImpl struct {
	Config     NBPConfig
	httpClient *http.Client
	auth       *authCache
}

func (c *NBPClientImpl) Hello() (*HelloResponse, error) {
	url := c.Config.BaseUrl + "/Hello"
	resp, err := c.httpClient.Get(url)

	if err != nil {
		//Unlikely
		return nil, err
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		//Unlikely
		return nil, err
	}

	raw := string(body)
	ret := &HelloResponse{
		StatusCode:  resp.StatusCode,
		RawResponse: raw,
		Data:        raw,
	}

	return ret, nil
}

func (c *NBPClientImpl) authenticate() (*authenticateResponse, error) {
	url := fmt.Sprintf("%s/Authenticate?Agency_Code=%s", c.Config.BaseUrl, c.Config.AgencyCode)
	basicAuth := fmt.Sprintf("%s:%s", c.Config.Username, c.Config.Password)
	basicAuth = base64.StdEncoding.EncodeToString([]byte(basicAuth))
	basicAuth = fmt.Sprintf("Basic %v", basicAuth)

	r, err := http.NewRequest("POST", url, nil)
	if err != nil {
		//Unlikely
		return nil, err
	}

	r.Header.Add("Authorization", basicAuth)

	resp, err := c.httpClient.Do(r)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 && resp.StatusCode != 400 {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			//Unlikely
			return nil, err
		}
		raw := string(body)
		return &authenticateResponse{
			StatusCode:  resp.StatusCode,
			RawResponse: raw,
		}, nil
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		//Unlikely
		return nil, err
	}

	auth := &authenticate{}
	derr := json.NewDecoder(bytes.NewBuffer(body)).Decode(auth)
	if derr != nil {
		//Fatal: Decode json should always success, need Alert
		return nil, derr
	}

	return &authenticateResponse{
		StatusCode:  resp.StatusCode,
		RawResponse: string(body),
	}, nil
}

func (c *NBPClientImpl) updateToken() {
	//Mutex
}

func (c *NBPClientImpl) BankList() (*BankListResponse, error) {
	return nil, nil
}

func (c *NBPClientImpl) AccountEnquiry(r AccountEnquiryRequest) (*AccountEnquiryResponse, error) {
	return nil, nil
}

func (c *NBPClientImpl) LoadRemittance(r LoadRemittanceRequest) (*LoadRemittanceResponse, error) {
	return nil, nil
}

func (c *NBPClientImpl) TransactionStatusByIds(r TransactionStatusByIdsRequest) (*TransactionStatusByIdsResponse, error) {
	return nil, nil
}

func (c *NBPClientImpl) TransactionStatusByDate(r TransactionStatusByDateRequest) (*TransactionStatusByDateResponse, error) {
	return nil, nil
}
func (c *NBPClientImpl) CancelTransaction(r CancelTransactionRequest) (*CancelTransactionResponse, error) {
	return nil, nil
}
