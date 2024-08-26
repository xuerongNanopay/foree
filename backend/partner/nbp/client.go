package nbp

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

type tokenData struct {
	token          string
	rawTokenExpiry string
	tokenExpiry    time.Time
}

type NBPClient interface {
	Hello() (*HelloResponse, error)
	BankList() (*BankListResponse, error)
	AccountEnquiry(AccountEnquiryRequest) (*AccountEnquiryResponse, error)
	LoadRemittanceCash(LoadRemittanceRequest) (*LoadRemittanceResponse, error)
	LoadRemittanceAccounts(LoadRemittanceRequest) (*LoadRemittanceResponse, error)
	LoadRemittanceThirdParty(LoadRemittanceRequest) (*LoadRemittanceResponse, error)
	TransactionStatusByIds(TransactionStatusByIdsRequest) (*TransactionStatusByIdsResponse, error)
	TransactionStatusByDate(TransactionStatusByDateRequest) (*TransactionStatusByDateResponse, error)
	CancelTransaction(CancelTransactionRequest) (*CancelTransactionResponse, error)
	GetConfigs() map[string]string
	SetConfig(key string, value string)
}

func NewNBPClient(configs map[string]string) NBPClient {
	nbpConfig := NewNBPConfigWithDefaultConfig(configs)

	return &NBPClientImpl{
		config: nbpConfig,
		httpClient: &http.Client{
			Timeout: 4 * time.Minute, // At least 3 minutes
		},
	}
}

type NBPClientImpl struct {
	config     NBPConfig
	httpClient *http.Client
	auth       tokenData
	mu         sync.Mutex
}

func (s *NBPClientImpl) GetConfigs() map[string]string {
	return s.config.ShowConfigs()
}

func (s *NBPClientImpl) SetConfig(key string, value string) {
	s.config.SetConfig(key, value)
}

func (c *NBPClientImpl) Hello() (*HelloResponse, error) {
	url := c.config.GetBaseUrl() + "/Hello"
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
		ResponseCommon: ResponseCommon{
			StatusCode:  resp.StatusCode,
			RawResponse: raw,
		},
	}

	return ret, nil
}

func (c *NBPClientImpl) authenticate() (*authenticateResponse, error) {
	url := fmt.Sprintf("%s/Authenticate?Agency_Code=%s", c.config.GetBaseUrl(), c.config.GetAgencyCode())
	basicAuth := fmt.Sprintf("%s:%s", c.config.GetAuthUsername(), c.config.GetAuthPassword())
	basicAuth = base64.StdEncoding.EncodeToString([]byte(basicAuth))
	basicAuth = fmt.Sprintf("Basic %v", basicAuth)

	r, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		//Unlikely; Fatal
		return nil, err
	}

	r.Header.Add("Authorization", basicAuth)

	resp, err := c.httpClient.Do(r)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		//Unlikely; Fatal
		return nil, err
	}

	auth := &authenticateResponse{
		ResponseCommon: ResponseCommon{
			StatusCode:  resp.StatusCode,
			RawResponse: string(body),
		},
	}

	if resp.StatusCode != 200 && resp.StatusCode != 400 {
		return auth, nil
	}

	derr := json.NewDecoder(bytes.NewBuffer(body)).Decode(auth)
	if derr != nil {
		//Fatal: Decode json should always success, need Alert
		return nil, derr
	}

	return auth, nil
}

// Update token but not throw error, Only log error
func (c *NBPClientImpl) maybeUpdateToken() {
	if isValidToken(c.auth, c.config.GetTokenExpiryThreshold()) {
		return
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	if isValidToken(c.auth, c.config.GetTokenExpiryThreshold()) {
		return
	}

	authResp, err := c.authenticate()
	if err != nil {
		//TODO: Fatal
		return
	}

	statusCode := authResp.StatusCode
	if statusCode/100 != 2 && statusCode/100 != 4 {
		//TODO: Error
		return
	}

	code := authResp.ResponseCode
	if code == "402" || code == "404" || code == "407" {
		//TODO: Fatal
		return
	}

	if code == "403" {
		//TODO: Error
		return
	}

	token := authResp.Token
	rawTokenExpiry := authResp.TokenExpiry
	tokenExpiry, err := parseTokenExpiryDate(rawTokenExpiry)
	if err != nil {
		//TODO: alarm/warming. We can't parse the time but we can still use the token.
		fmt.Printf("NBP Client authenticate: unable to parse token_expiry `%v`", rawTokenExpiry)
	}

	c.auth.token = token
	c.auth.rawTokenExpiry = rawTokenExpiry
	c.auth.tokenExpiry = tokenExpiry
}

func (c *NBPClientImpl) BankList() (*BankListResponse, error) {
	c.maybeUpdateToken()
	url := fmt.Sprintf("%s/BankList", c.config.GetBaseUrl())

	r := requestCommon{
		Token:      c.auth.token,
		AgencyCode: c.config.GetAgencyCode(),
	}

	rawReqeust, err := json.Marshal(r)
	if err != nil {
		//Unlikely; Fatal
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(rawReqeust))
	if err != nil {
		//Unlikely; Fatal
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
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

	ret := &BankListResponse{
		ResponseCommon: ResponseCommon{
			StatusCode:  resp.StatusCode,
			RawRequest:  string(rawReqeust),
			RawResponse: string(body),
		},
	}

	if resp.StatusCode != 200 && resp.StatusCode != 400 {
		return ret, nil
	}

	derr := json.NewDecoder(bytes.NewBuffer(body)).Decode(ret)
	if derr != nil {
		//Unlikely; Fatal
		return nil, err
	}
	return ret, nil
}

func (c *NBPClientImpl) AccountEnquiry(r AccountEnquiryRequest) (*AccountEnquiryResponse, error) {
	c.maybeUpdateToken()
	url := fmt.Sprintf("%s/AccountEnquiry", c.config.GetBaseUrl())

	r.Token = c.auth.token
	r.AgencyCode = c.config.GetAgencyCode()

	rawReqeust, err := json.Marshal(r)
	if err != nil {
		//Unlikely; Fatal
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(rawReqeust))
	if err != nil {
		//Unlikely; Fatal
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
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

	ret := &AccountEnquiryResponse{
		ResponseCommon: ResponseCommon{
			StatusCode:  resp.StatusCode,
			RawRequest:  string(rawReqeust),
			RawResponse: string(body),
		},
	}

	if resp.StatusCode != 200 && resp.StatusCode != 400 {
		return ret, nil
	}

	derr := json.NewDecoder(bytes.NewBuffer(body)).Decode(ret)
	if derr != nil {
		//Unlikely; Fatal
		return nil, err
	}
	return ret, nil

}

func (c *NBPClientImpl) LoadRemittanceCash(r LoadRemittanceRequest) (*LoadRemittanceResponse, error) {
	r.PmtMode = PMTModeCash
	return c.loadRemittance("LoadRemittanceCash", r)
}

func (c *NBPClientImpl) LoadRemittanceAccounts(r LoadRemittanceRequest) (*LoadRemittanceResponse, error) {
	r.PmtMode = PMTModeAccountTransfers
	return c.loadRemittance("LoadRemittanceAccounts", r)
}

func (c *NBPClientImpl) LoadRemittanceThirdParty(r LoadRemittanceRequest) (*LoadRemittanceResponse, error) {
	r.PmtMode = PMTModeThirdPartyPayments
	return c.loadRemittance("LoadRemittanceThirdParty", r)
}

func (c *NBPClientImpl) loadRemittance(endpoint string, r LoadRemittanceRequest) (*LoadRemittanceResponse, error) {
	c.maybeUpdateToken()
	url := fmt.Sprintf("%s/%s", c.config.GetBaseUrl(), endpoint)
	r.Token = c.auth.token
	r.AgencyCode = c.config.GetAgencyCode()

	rawReqeust, err := json.Marshal(r)
	if err != nil {
		//Unlikely; Fatal
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(rawReqeust))
	if err != nil {
		//Unlikely; Fatal
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
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

	ret := &LoadRemittanceResponse{
		ResponseCommon: ResponseCommon{
			StatusCode:  resp.StatusCode,
			RawRequest:  string(rawReqeust),
			RawResponse: string(body),
		},
	}

	if resp.StatusCode != 200 && resp.StatusCode != 400 {
		return ret, nil
	}

	derr := json.NewDecoder(bytes.NewBuffer(body)).Decode(ret)
	if derr != nil {
		//Unlikely; Fatal
		return nil, err
	}
	return ret, nil
}

func (c *NBPClientImpl) TransactionStatusByIds(r TransactionStatusByIdsRequest) (*TransactionStatusByIdsResponse, error) {
	c.maybeUpdateToken()
	url := fmt.Sprintf("%s/TransactionStatusByIds", c.config.GetBaseUrl())

	r.Token = c.auth.token
	r.AgencyCode = c.config.GetAgencyCode()

	rawReqeust, err := json.Marshal(r)
	if err != nil {
		//Unlikely; Fatal
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(rawReqeust))
	if err != nil {
		//Unlikely; Fatal
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
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

	ret := &TransactionStatusByIdsResponse{
		ResponseCommon: ResponseCommon{
			StatusCode:  resp.StatusCode,
			RawRequest:  string(rawReqeust),
			RawResponse: string(body),
		},
	}

	if resp.StatusCode != 200 && resp.StatusCode != 400 {
		return ret, nil
	}

	derr := json.NewDecoder(bytes.NewBuffer(body)).Decode(ret)
	if derr != nil {
		//Unlikely; Fatal
		return nil, err
	}
	return ret, nil
}

func (c *NBPClientImpl) TransactionStatusByDate(r TransactionStatusByDateRequest) (*TransactionStatusByDateResponse, error) {
	c.maybeUpdateToken()
	url := fmt.Sprintf("%s/TransactionStatus", c.config.GetBaseUrl())

	r.Token = c.auth.token
	r.AgencyCode = c.config.GetAgencyCode()

	rawReqeust, err := json.Marshal(r)
	if err != nil {
		//Unlikely; Fatal
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(rawReqeust))
	if err != nil {
		//Unlikely; Fatal
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
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

	ret := &TransactionStatusByDateResponse{
		ResponseCommon: ResponseCommon{
			StatusCode:  resp.StatusCode,
			RawRequest:  string(rawReqeust),
			RawResponse: string(body),
		},
	}

	if resp.StatusCode != 200 && resp.StatusCode != 400 {
		return ret, nil
	}

	derr := json.NewDecoder(bytes.NewBuffer(body)).Decode(ret)
	if derr != nil {
		//Unlikely; Fatal
		return nil, err
	}
	return ret, nil
}

func (c *NBPClientImpl) CancelTransaction(r CancelTransactionRequest) (*CancelTransactionResponse, error) {
	c.maybeUpdateToken()
	url := fmt.Sprintf("%s/CancelTransaction", c.config.GetBaseUrl())

	r.Token = c.auth.token
	r.AgencyCode = c.config.GetAgencyCode()

	rawReqeust, err := json.Marshal(r)
	if err != nil {
		//Unlikely; Fatal
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(rawReqeust))
	if err != nil {
		//Unlikely; Fatal
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
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

	ret := &CancelTransactionResponse{
		ResponseCommon: ResponseCommon{
			StatusCode:  resp.StatusCode,
			RawRequest:  string(rawReqeust),
			RawResponse: string(body),
		},
	}

	if resp.StatusCode != 200 && resp.StatusCode != 400 {
		return ret, nil
	}

	derr := json.NewDecoder(bytes.NewBuffer(body)).Decode(ret)
	if derr != nil {
		//Unlikely; Fatal
		return nil, err
	}
	return ret, nil
}
