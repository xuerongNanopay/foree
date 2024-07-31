package nbp

import (
	"fmt"
	"strings"
	"time"
)

const (
	TxStatusRejected  = "REJECTED"
	TxStatusError     = "ERROR"
	TxStatusPaid      = "PAID"
	TxStatusCancelled = "CANCELLED"
	TxStatusInProcess = "IN_PROCESS"
	TxStatusPending   = "PENDING"
)

type responseGetter interface {
	GetStatusCode() int
	GetRawRequest() string
	GetRawResponse() string
	GetResponseCode() string
	GetResponseMessage() string
}

type ResponseCommon struct {
	StatusCode      int    `json:"-"`
	RawRequest      string `json:"-"`
	RawResponse     string `json:"-"`
	ResponseCode    string `json:"ResponseCode"`
	ResponseMessage string `json:"ResponseMessage"`
}

func (r *ResponseCommon) GetStatusCode() int {
	return r.StatusCode
}

func (r *ResponseCommon) GetRawRequest() string {
	return r.RawRequest
}

func (r *ResponseCommon) GetRawResponse() string {
	return r.RawResponse
}

func (r *ResponseCommon) GetResponseCode() string {
	return r.ResponseCode
}

func (r *ResponseCommon) GetResponseMessage() string {
	return r.ResponseMessage
}

type HelloResponse struct {
	ResponseCommon
}

type authenticateResponse struct {
	ResponseCommon
	Token       string `json:"Token"`
	TokenExpiry string `json:"Token_Expiry"`
}

type BankListResponse struct {
	ResponseCommon
	Banklist []BankListEntry `json:"banklist"`
}

type BankListEntry struct {
	BankName string `json:"bankName"`
	BankAbbr string `json:"bankAbbr"`
}

type AccountEnquiryResponse struct {
	ResponseCommon
	IBAN          string `json:"IBAN"`
	AccountNo     string `json:"AccountNo"`
	AccountTitle  string `json:"AccountTitle"`
	BranchCode    int32  `json:"BranchCode"`
	AccountStatus string `json:"AccountStatus"`
	BankName      string `json:"BankName"`
}

type LoadRemittanceResponse struct {
	ResponseCommon
	GlobalId   string `json:"Global_Id"`
	TrackingId string `json:"Tracking_Id"`
}

type CancelTransactionResponse struct {
	ResponseCommon
	GlobalId string `json:"Global_Id"`
}

type TransactionStatus struct {
	GlobalId                string `json:"Global_Id"`
	TrackingId              string `json:"Tracking_Id"`
	Status                  string `json:"status"`
	StatusDetails           string `json:"Status_Details"`
	BeneficiaryIdType       string `json:"Beneficiary_Id_Type"`
	BeneficiaryIdNumber     string `json:"Beneficiary_Id_Number"`
	BranchCode              uint32 `json:"Branch_Code"`
	BranchName              string `json:"Branch_Name"`
	BeneficiaryName         string `json:"Beneficiary_Name"`
	BeneficiaryIdExpiryDate string `json:"Beneficiary_Id_Expiry_Date"`
	BeneficiaryIdIssueDate  string `json:"Beneficiary_Id_Issue_Date"`
	BeneficiaryIdIssuedBy   string `json:"Beneficiary_Id_Issued_By"`
	BeneficiaryDateOfBirth  string `json:"Beneficiary_Date_Of_Birth"`
	ProcessingDate          string `json:"Processing_Date"`
	ProcessingTime          string `json:"Processing_Time"`
}

type TransactionStatusByIdsResponse struct {
	ResponseCommon
	TransactionStatuses []TransactionStatus `json:"transactionStatuses"`
}

type TransactionStatusByDateResponse struct {
	ResponseCommon
	TransactionStatuses []TransactionStatus `json:"transactionStatuses"`
}

type requestCommon struct {
	Token      string `json:"Token,omitempty"`
	AgencyCode string `json:"Agency_Code,omitempty"`
}

type BankListRequest struct {
	requestCommon
}

type AccountEnquiryRequest struct {
	requestCommon
	AccountNo  string `json:"AccountNo,omitempty"`
	BranchCode int32  `json:"BranchCode,omitempty"`
	BankName   string `json:"BankName,omitempty"`
}

type NBPAmount float64

func (a NBPAmount) MarshalJSON() ([]byte, error) {
	s := fmt.Sprintf("%.2f", a)
	return []byte(s), nil
}

type NBPDate time.Time

func (d NBPDate) MarshalJSON() ([]byte, error) {
	t := time.Time(d)
	s := t.Format(time.DateOnly)
	return []byte("\"" + s + "\""), nil
}

type LoadRemittanceRequest struct {
	requestCommon
	GlobalId                        string    `json:"Global_Id,omitempty"`
	Currency                        string    `json:"Currency,omitempty"`
	Amount                          NBPAmount `json:"Amount,omitempty"` //see: https://stackoverflow.com/questions/61811463/golang-encode-float-to-json-with-specified-precision
	PmtMode                         string    `json:"Pmt_Mode,omitempty"`
	RemitterName                    string    `json:"Remitter_Name,omitempty"`
	RemitterAddress                 string    `json:"Remitter_Address,omitempty"`
	RemitterEmail                   string    `json:"Remitter_Email,omitempty"`
	RemitterContact                 string    `json:"Remitter_Contact,omitempty"`
	RemitterIdType                  string    `json:"Remitter_Id_Type,omitempty"`
	RemitterId                      string    `json:"Remitter_Id,omitempty"`
	BeneficiaryName                 string    `json:"Beneficiary_Name,omitempty"`
	BeneficiaryAddress              string    `json:"Beneficiary_Address,omitempty"`
	BeneficiaryContact              string    `json:"Beneficiary_Contact,omitempty"`
	BeneficiaryExpectedId           string    `json:"Beneficiary_Expectedid,omitempty"`
	BeneficiaryBank                 string    `json:"Beneficiary_Bank,omitempty"`
	BeneficiaryBranch               string    `json:"Beneficiary_Branch,omitempty"`
	BeneficiaryAccount              string    `json:"Beneficiary_Account,omitempty"`
	PurposeRemittance               string    `json:"Purpose_Remittance,omitempty"`
	BeneficiaryCity                 string    `json:"Beneficiary_City,omitempty"`
	OriginatingCountry              string    `json:"Originating_Country,omitempty"`
	TransactionDate                 string    `json:"Transaction_Date,omitempty"` //yyyy-MM-dd
	RemitterAccountNo               string    `json:"remitter_AccountNo,omitempty"`
	RemitterFatherName              string    `json:"remitter_FatherName,omitempty"`
	RemitterDOB                     *NBPDate  `json:"remitter_DOB,omitempty"` //yyyy-MM-dd
	RemitterPOB                     *NBPDate  `json:"remitter_POB,omitempty"`
	RemitterNationality             string    `json:"remitter_Nationality,omitempty"`
	RemitterBeneficiaryRelationship string    `json:"remitter_BeneficiaryRelationship,omitempty"`
}

type NBPIds []string

func (d NBPIds) MarshalJSON() ([]byte, error) {
	return []byte("\"" + strings.Join(d, ",") + "\""), nil
}

type TransactionStatusByIdsRequest struct {
	requestCommon
	Ids NBPIds `json:"Ids,omitempty"`
}

type TransactionStatusByDateRequest struct {
	requestCommon
	Date *NBPDate `json:"Date,omitempty"`
}

type CancelTransactionRequest struct {
	requestCommon
	GlobalId           string `json:"Global_Id,omitempty"`
	CancellationReason string `json:"Cancellation_Reason,omitempty"`
}
