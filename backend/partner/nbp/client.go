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

type authCache struct {
	token          string
	rawTokenExpiry string
	tokenExpiry    *time.Time
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
	mu         sync.Mutex
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
		ResponseCommon: ResponseCommon{
			StatusCode:  resp.StatusCode,
			RawResponse: raw,
		},
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
			ResponseCommon: ResponseCommon{
				StatusCode:  resp.StatusCode,
				RawResponse: raw,
			},
		}, nil
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		//Unlikely
		return nil, err
	}

	auth := &authenticateResponse{
		ResponseCommon: ResponseCommon{
			StatusCode:  resp.StatusCode,
			RawResponse: string(body),
		},
	}
	derr := json.NewDecoder(bytes.NewBuffer(body)).Decode(auth)
	if derr != nil {
		//Fatal: Decode json should always success, need Alert
		return nil, derr
	}

	return auth, nil
}

func (c *NBPClientImpl) updateToken() error {
	if isTokenAvailable(c.auth, c.Config.TokenExpiryThreshold) {
		return nil
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	if isTokenAvailable(c.auth, c.Config.TokenExpiryThreshold) {
		return nil
	}

	authResp, err := c.authenticate()
	if err != nil {
		//TODO: Fatal
		return fmt.Errorf("NBP Client authenticate: raise error `%s`", err.Error())
	}

	statusCode := authResp.StatusCode
	if statusCode != 200 && statusCode != 400 {
		//TODO: Error
		return fmt.Errorf("NBP Client authenticate: status code `%v` response body `%s`", authResp.StatusCode, authResp.RawResponse)
	}

	code := authResp.ResponseCode
	if code == "402" || code == "404" || code == "407" {
		//TODO: Fatal
		return fmt.Errorf("NBP Client authenticate: status code `%v` response body `%s`", authResp.StatusCode, authResp.RawResponse)
	}

	if code == "403" {
		//TODO: Error
		return fmt.Errorf("NBP Client authenticate: status code `%v` response body `%s`", authResp.StatusCode, authResp.RawResponse)
	}

	token := authResp.Token
	rawTokenExpiry := authResp.TokenExpiry
	tokenExpiry, err := parseTokenExpiryDate(rawTokenExpiry)
	if err != nil {
		//TODO: alarm/warming. We can't parse the time but we can still use the token.
		fmt.Printf("NBP Client authenticate: unable to parse token_expiry `%v`", rawTokenExpiry)
	}

	auth := &authCache{
		token:          token,
		rawTokenExpiry: rawTokenExpiry,
		tokenExpiry:    tokenExpiry,
	}

	c.auth = auth

	return nil
}

func (c *NBPClientImpl) BankList() (*BankListResponse, error) {
	return nil, nil
}

func (c *NBPClientImpl) AccountEnquiry(r AccountEnquiryRequest) (*AccountEnquiryResponse, error) {
	resp, err := c.retry(func() (responseGetter, error) {
		ret := &AccountEnquiryResponse{}
		return ret, nil
	})
	if err != nil {
		return nil, err
	}
	return resp.(*AccountEnquiryResponse), nil
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

func (c *NBPClientImpl) retry(f func() (responseGetter, error)) (responseGetter, error) {
	attempts := c.Config.AuthAttempts
	if attempts < 1 {
		attempts = 3
	}

	var tokenErr error
	for i := 0; i < attempts; i++ {
		tokenErr = c.updateToken()
		if tokenErr != nil {
			time.Sleep(4 * time.Second)
		} else {
			r, err := f()
			if err != nil {
				return r, err
			}

			if r.GetResponseCode() == "401" {
				//TODO: log
				//Reset authCache if the token is already expired.
				c.mu.Lock()
				c.auth = nil
				c.mu.Unlock()
				time.Sleep(4 * time.Second)
			} else {
				return r, nil
			}

		}
	}
	return nil, tokenErr
}
