package scotia

import (
	"bytes"
	cryptoRsa "crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	tokenExpiryInMinutes      = 5
	tokenExpiryThresholdInSec = 30
)

type ScotiaClient interface {
	RequestPayment(req RequestPaymentRequest) (*RequestPaymentResponse, error)

	GetConfigs() map[string]string
	SetConfig(key string, value string)
}

type tokenData struct {
	token       string
	tokenExpiry time.Time
	scope       string
	tokenType   string
	expirysIn   int64
}

type rsa struct {
	privateKeyDir string
	signKey       *cryptoRsa.PrivateKey
}

func initRSA(config ScotiaConfig) (*rsa, error) {
	//TODO: load rsa from config.
	return nil, nil
}

func NewScotiaClientImpl(configs map[string]string) ScotiaClient {
	scotiaConfig := NewScotiaConfigWithDefaultConfig(configs)

	return &scotiaClientImpl{
		config: scotiaConfig,
	}
}

type scotiaClientImpl struct {
	config     ScotiaConfig
	rsa        *rsa
	auth       *tokenData
	mu         sync.Mutex
	httpClient *http.Client
}

func (s *scotiaClientImpl) GetConfigs() map[string]string {
	return s.config.ShowConfigs()
}

func (s *scotiaClientImpl) SetConfig(key string, value string) {
	s.config.SetConfig(key, value)
}

func (s *scotiaClientImpl) RequestPayment(req RequestPaymentRequest) (*RequestPaymentResponse, error) {
	url := fmt.Sprintf("%s/treasury/payments/rtp/v1/requests", s.config.GetBaseUrl())

	token, err := s.getToken()
	if err != nil {
		return nil, err
	}

	basicAuth := fmt.Sprintf("Basic %v", token)

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

	s.setCommonHeaders(r)
	r.Header.Add("Content-Type", "application/json")
	r.Header.Add("Authorization", basicAuth)

	resp, err := s.httpClient.Do(r)
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

	ret := &RequestPaymentResponse{
		ResponseCommon: ResponseCommon{
			StatusCode:  resp.StatusCode,
			RawRequest:  string(rawReqeust),
			RawResponse: string(body),
		},
	}

	derr := json.NewDecoder(bytes.NewBuffer(body)).Decode(ret)
	if derr != nil {
		//TODO: Logger error. return token caller should hanlde the Error
		return ret, nil
	}

	return ret, nil

}

func (s *scotiaClientImpl) tokenRequest() (*tokenResponse, error) {
	endpoint := fmt.Sprintf("%s/scotiabank/wam/vi/getToken", s.config.GetBaseUrl())
	basicAuth := fmt.Sprintf("%s:%s", s.config.GetAuthUsername(), s.config.GetAuthPassword())
	basicAuth = base64.StdEncoding.EncodeToString([]byte(basicAuth))
	basicAuth = fmt.Sprintf("Basic %v", basicAuth)

	jwt, err := s.signJWT()
	if err != nil {
		return nil, err
	}

	formData := url.Values{}
	formData.Add("grant_type", "client_credentials")
	formData.Add("scope", s.config.GetScope())
	formData.Add("client_id", s.config.GetClientId())
	formData.Add("client_assertion", jwt)
	formData.Add("client_assertion_type", "urn:ietf:params:oauth:client-assertion-type:jwt-bearer")

	r, err := http.NewRequest(http.MethodPost, endpoint, strings.NewReader(formData.Encode()))
	if err != nil {
		//Unlikely; Fatal
		return nil, err
	}

	r.Header.Add("Authorization", basicAuth)
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := s.httpClient.Do(r)
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

	token := &tokenResponse{
		ResponseCommon: ResponseCommon{
			StatusCode:  resp.StatusCode,
			RawResponse: string(body),
		},
	}

	derr := json.NewDecoder(bytes.NewBuffer(body)).Decode(token)
	if derr != nil {
		//TODO: Logger error. return token caller should hanlde the Error
		return token, nil
	}

	return token, nil

}

func (s *scotiaClientImpl) signJWT() (string, error) {
	claims := &jwt.RegisteredClaims{
		Subject:   s.config.GetClientId(),
		Audience:  []string{s.config.GetJWTAudience()},
		Issuer:    s.config.GetClientId(),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(s.config.GetJWTExpiryMinutes()) * time.Minute)),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	token.Header["kid"] = s.config.GetJWTKid()
	ss, err := token.SignedString(s.rsa.signKey)
	if err != nil {
		return "", fmt.Errorf("ScotiaClientImpl JWT signature got error `%v`", err.Error())
	}
	return ss, nil
}

func (s *scotiaClientImpl) setCommonHeaders(r *http.Request) {
	r.Header.Add("customer-profile-id", s.config.GetProfileId())
	r.Header.Add("x-country-code", s.config.GetCountryCode())
	r.Header.Add("x-api-key", s.config.GetApiKey())
}

// If token is invalid, then update token
func (s *scotiaClientImpl) getToken() (string, error) {
	err := s.maybeUpdateToken()
	if err != nil {
		return "", err
	}
	return s.auth.token, nil
}

func (s *scotiaClientImpl) maybeUpdateToken() error {
	if isValidToken(s.auth, tokenExpiryThresholdInSec) {
		return nil
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if isValidToken(s.auth, tokenExpiryThresholdInSec) {
		return nil
	}
	tokenResp, err := s.tokenRequest()
	if err != nil {
		return err
	}

	if tokenResp.StatusCode/100 != 2 {
		return fmt.Errorf("scotialClientImpl: token request failed with status code `%v`, response payload `%v`", tokenResp.StatusCode, tokenResp.RawResponse)
	}

	newAuth := &tokenData{
		token:       tokenResp.AccessToken,
		tokenType:   tokenResp.TokenType,
		scope:       tokenResp.Scope,
		expirysIn:   tokenResp.ExpiresIn,
		tokenExpiry: time.Now().Add(tokenExpiryInMinutes * time.Minute),
	}

	s.auth = newAuth

	return nil
}

func isValidToken(auth *tokenData, threshold int64) bool {
	if auth == nil || auth.token == "" || auth.tokenExpiry.IsZero() {
		return false
	}

	if time.Now().Unix()+threshold >= auth.tokenExpiry.Unix() {
		return false
	}

	return true
}
