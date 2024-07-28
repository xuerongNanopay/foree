package nbp

type responseWrapper[R any] struct {
	StatusCode  int
	RawRequest  string
	RawResponse string
	Data        R
}

type authenticate struct {
	ResponseCode    string `json:"ResponseCode"`
	ResponseMessage string `json:"ResponseMessage"`
	Token           string `json:"Token"`
	TokenExpiry     string `json:"Token_Expiry"`
}

type BankList struct {
	ResponseCode    string          `json:"ResponseCode"`
	ResponseMessage string          `json:"ResponseMessage"`
	Banklist        []BankListEntry `json:"banklist"`
}

type BankListEntry struct {
	BankName string `json:"bankName"`
	BankAbbr string `json:"bankAbbr"`
}

type AccountEnquiry struct {
	ResponseCode    string `json:"ResponseCode"`
	ResponseMessage string `json:"ResponseMessage"`
	IBAN            string `json:"IBAN"`
	AccountNo       string `json:"AccountNo"`
	AccountTitle    string `json:"AccountTitle"`
	BranchCode      int32  `json:"BranchCode"`
	AccountStatus   string `json:"AccountStatus"`
	BankName        string `json:"BankName"`
}

type LoadRemittance struct {
	ResponseCode    string `json:"ResponseCode"`
	ResponseMessage string `json:"ResponseMessage"`
	GlobalId        string `json:"Global_Id"`
	TrackingId      string `json:"Tracking_Id"`
}

type CancelTransaction struct {
	ResponseCode    string `json:"ResponseCode"`
	ResponseMessage string `json:"ResponseMessage"`
	GlobalId        string `json:"Global_Id"`
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

type TransactionStatuses struct {
	ResponseCode        string              `json:"ResponseCode"`
	ResponseMessage     string              `json:"ResponseMessage"`
	TransactionStatuses []TransactionStatus `json:"transactionStatuses"`
}

type HelloResponse responseWrapper[string]
type authenticateResponse responseWrapper[authenticate]
type BankListResponse responseWrapper[BankList]
type AccountEnquiryResponse responseWrapper[AccountEnquiry]
type LoadRemittanceResponse responseWrapper[LoadRemittance]
type TransactionStatusByIdsResponse responseWrapper[TransactionStatuses]
type TransactionStatusByDateResponse responseWrapper[TransactionStatuses]
type CancelTransactionResponse responseWrapper[CancelTransaction]

type AuthenticateRequest struct {
	AgencyCode int32
	UserName   string
	Password   string
}

type BankListRequest struct {
	token      string `json:"Token"`
	agencyCode string `json:"Agency_Code"`
}

type AccountEnquiryRequest struct {
	token      string `json:"Token"`
	agencyCode string `json:"Agency_Code"`
	AccountNo  string `json:"AccountNo"`
	BranchCode int32  `json:"AccountNo"`
	BankName   string `json:"BankName"`
}

type LoadRemittanceRequest struct {
	token                           string  `json:"Token"`
	agencyCode                      string  `json:"Agency_Code"`
	GlobalId                        string  `json:"Global_Id"`
	Currency                        string  `json:"Currency"`
	Amount                          float64 `json:"Amount"` //see: https://stackoverflow.com/questions/61811463/golang-encode-float-to-json-with-specified-precision
	PmtMode                         string  `json:"Pmt_Mode"`
	RemitterName                    string  `json:"Remitter_Name"`
	RemitterAddress                 string  `json:"Remitter_Address"`
	RemitterEmail                   string  `json:"Remitter_Email"`
	RemitterContact                 string  `json:"Remitter_Contact"`
	RemitterIdType                  string  `json:"Remitter_Id_Type"`
	RemitterId                      string  `json:"Remitter_Id"`
	BeneficiaryName                 string  `json:"Beneficiary_Name"`
	BeneficiaryAddress              string  `json:"Beneficiary_Address"`
	BeneficiaryContact              string  `json:"Beneficiary_Contact"`
	BeneficiaryExpectedId           string  `json:"Beneficiary_Expectedid"`
	BeneficiaryBank                 string  `json:"Beneficiary_Bank"`
	BeneficiaryBranch               string  `json:"Beneficiary_Branch"`
	BeneficiaryAccount              string  `json:"Beneficiary_Account"`
	PurposeRemittance               string  `json:"Purpose_Remittance"`
	BeneficiaryCity                 string  `json:"Beneficiary_City"`
	OriginatingCountry              string  `json:"Originating_Country"`
	TransactionDate                 string  `json:"Transaction_Date"` //yyyy-MM-dd
	RemitterAccountNo               string  `json:"remitter_AccountNo"`
	RemitterFatherName              string  `json:"remitter_FatherName"`
	RemitterDOB                     string  `json:"remitter_DOB"` //yyyy-MM-dd
	RemitterPOB                     string  `json:"remitter_POB"`
	RemitterNationality             string  `json:"remitter_Nationality"`
	RemitterBeneficiaryRelationship string  `json:"remitter_BeneficiaryRelationship"`
}

type TransactionStatusByIdsRequest struct {
	token      string `json:"Token"`
	agencyCode string `json:"Agency_Code"`
	Ids        string `json:"Ids"`
}

type TransactionStatusByDateRequest struct {
	token      string `json:"Token"`
	agencyCode string `json:"Agency_Code"`
	Date       string `json:"Date"`
}

type CancelTransactionRequest struct {
	token              string `json:"Token"`
	agencyCode         string `json:"Agency_Code"`
	GlobalId           string `json:"Global_Id"`
	CancellationReason string `json:"Cancellation_Reason"`
}
