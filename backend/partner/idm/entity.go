package idm

import (
	"fmt"
	"time"
)

const (
	ResultStatusAccept        = "ACCEPT"
	ResultStatusDeny          = "DENY"
	ResultStatusMannualReview = "MANUAL_REVIEW"
)

type IDMAmount float64

func (a IDMAmount) MarshalJSON() ([]byte, error) {
	s := fmt.Sprintf("%.2f", a)
	return []byte(s), nil
}

type IDMDate time.Time

func (d IDMDate) MarshalJSON() ([]byte, error) {
	t := time.Time(d)
	s := t.Format(time.DateOnly)
	return []byte("\"" + s + "\""), nil
}

type IDMRequest struct {
	UserAccountIdentifier   string    `json:"man,omitempty"`
	UserEmail               string    `json:"tea,omitempty"`
	Ip                      string    `json:"ip,omitempty"`
	PhoneNumber             string    `json:"phn,omitempty"`
	Dob                     *IDMDate  `json:"dob,omitempty"`
	BillingFirstName        string    `json:"bfn,omitempty"`
	BillingMiddleName       string    `json:"bmn,omitempty"`
	BillingLastname         string    `json:"bln,omitempty"`
	BillingStreet           string    `json:"bsn,omitempty"`
	BillingCity             string    `json:"bc,omitempty"`
	BillingState            string    `json:"bs,omitempty"`
	BillingPostcode         string    `json:"bz,omitempty"`
	BillingCountry          string    `json:"bco,omitempty"`
	ShippingFirstName       string    `json:"sfn,omitempty"`
	ShippingMiddleName      string    `json:"smn,omitempty"`
	ShippingLastname        string    `json:"sln,omitempty"`
	ShippingStreet          string    `json:"ssn,omitempty"`
	ShippingCity            string    `json:"sc,omitempty"`
	ShippingState           string    `json:"ss,omitempty"`
	ShippingPostcode        string    `json:"sz,omitempty"`
	ShippingCountry         string    `json:"sco,omitempty"`
	DestPhoneNumber         string    `json:"dph,omitempty"`
	DestACHHash             string    `json:"dpach,omitempty"`
	DestDigitalAccNOHash    string    `json:"dphash,omitempty"`
	SrcACHHash              string    `json:"pach,omitempty"`
	SrcDigitalAccNOHash     string    `json:"phash,omitempty"`
	TransactionIdentifier   string    `json:"tid,omitempty"`
	TransactionCreationTime int64     `json:"tti,omitempty"`
	Amount                  IDMAmount `json:"amt,omitempty"`
	Currency                string    `json:"ccy,omitempty"`
	Profile                 string    `json:"profile,omitempty"`
	Tags                    []string  `json:"tags,omitempty"`
	Nationality             string    `json:"nationality,omitempty"`
	NationalId              string    `json:"nationalId,omitempty"`
	TaxId                   string    `json:"taxId,omitempty"`
	VoterId                 string    `json:"voterId,omitempty"`
	DriverId                string    `json:"driverId,omitempty"`
	PassportId              string    `json:"passportId,omitempty"`
	BeneBankName            string    `json:"memo9,omitempty"`
	PayoutCurrency          string    `json:"memo12,omitempty"`
	PayoutAmount            IDMAmount `json:"memo13,omitempty"`
	TransactionRefId        string    `json:"memo14,omitempty"`
	IsCashPickup            bool      `json:"memo17,omitempty"`
	RemitterPOB             string    `json:"memo18,omitempty"`
	RemitterCOR             string    `json:"memo19,omitempty"` // sender country of residence
	SRRelationship          string    `json:"memo20,omitempty"` // sender and receiver relationship
	PurposeOfTransfer       string    `json:"memo15,omitempty"`
	RemitterOccupation      string    `json:"memo16,omitempty"`
}

type ResponseCommon struct {
	StatusCode  int    `json:"-"`
	RawRequest  string `json:"-"`
	RawResponse string `json:"-"`
}

type IDMResponse struct {
	ResponseCommon
	CurrentUserReputation           string `json:"user"`
	PreviousUserReputation          string `json:"upr"`
	FraudEvaluationResult           string `json:"frp"`
	FraudRuleName                   string `json:"frn"`
	FraudRuleDescription            string `json:"frd"`
	TransactionId                   string `json:"tid"`
	UserReputationDescription       string `json:"erd"`
	PolicyEvalutionResult           string `json:"res"`
	TransactionEvaluationResultCode []int  `json:"rcd"`
}

func (r *IDMResponse) GetResultStatus() string {
	if r.PolicyEvalutionResult == "" {
		return r.FraudEvaluationResult
	}
	return r.PolicyEvalutionResult
}
