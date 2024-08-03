package scotia

import (
	"fmt"
	"time"
)

const (
	ProductionCodeDomestic = "DOMESTIC"
	CDIndicatorCREDIT      = "CRDT"
	CDIndicatorDebit       = "DBIT"
	LanguageEN             = "EN"
	LanguageFR             = "FR"
	CurrencyCAD            = "CAD"
)

type ScotiaDatetime time.Time

func (d ScotiaDatetime) MarshalJSON() ([]byte, error) {
	t := time.Time(d)
	s := t.Format(time.RFC3339)
	return []byte("\"" + s + "\""), nil
}

type ScotiaAmount float64

func (a ScotiaAmount) MarshalJSON() ([]byte, error) {
	s := fmt.Sprintf("%.2f", a)
	return []byte(s), nil
}

type ScotiaAmtData struct {
	Amount   ScotiaAmount `json:"amount,omitempty"`
	Currency string       `json:"currency,omitempty"`
}

type SchemeNameData struct {
	Code        string `json:"code,omitempty"`
	Proprietary string `json:"proprietary,omitempty"`
}

type OtherData struct {
	Identification string          `json:"identification,omitempty"`
	SchemeName     *SchemeNameData `json:"scheme_name,omitempty"`
}

type OrganisationIdentificationData struct {
	Other []OtherData `json:"other,omitempty"`
}

type IdentificationData struct {
	OrganisationIdentification *OrganisationIdentificationData `json:"organisation_identification,omitempty"`
}

type InitiatingPartyData struct {
	Name               string              `json:"name,omitempty"`
	Identification     *IdentificationData `json:"identification,omitempty"`
	CountryOfResidence string              `json:"country_of_residence,omitempty"`
}

type ContactDetailsData struct {
	EmailAddress string `json:"email_address,omitempty"`
}

type DebtorData struct {
	Name               string              `json:"name,omitempty"`
	CountryOfResidence string              `json:"country_of_residence,omitempty"`
	ContactDetails     *ContactDetailsData `json:"contact_details,omitempty"`
}

type CreditorData struct {
	Name               string              `json:"name,omitempty"`
	CountryOfResidence string              `json:"country_of_residence,omitempty"`
	ContactDetails     *ContactDetailsData `json:"contact_details,omitempty"`
}

type CreditorAccountData struct {
	Identification string `json:"identification,omitempty"`
	Currency       string `json:"currency,omitempty"`
	SchemeName     string `json:"scheme_nam,omitempty"`
}

type FraudSupplementaryInfoData struct {
	CustomerAuthenticationMethod string `json:"customer_authentication_method,omitempty"`
	CustomerIpAddress            string `json:"customer_ip_address,omitempty"`
}

type PaymentConditionData struct {
	AmountModificationAllowed  bool `json:"amount_modification_allowed,omitempty"`
	EarlyPaymentAllowed        bool `json:"early_payment_allowed,omitempty"`
	GuaranteedPaymentRequested bool `json:"guaranteed_payment_requested,omitempty"`
}

type CategoryPurposeData struct {
	Code string `json:"code,omitempty"`
}

type PaymentTypeInformationData struct {
	CategoryPurpose *CategoryPurposeData `json:"category_purpose,omitempty"`
}

type RemittanceInformationData struct {
	Unstructured []string `json:"unstructured,omitempty"`
}

type RequestPaymentRequestData struct {
	ProductCode                    string                      `json:"product_code,omitempty"`
	MessageIdentification          string                      `json:"message_identification,omitempty"`
	EndToEndIdentification         string                      `json:"end_to_end_identification,omitempty"`
	CreditDebitIndicator           string                      `json:"credit_debit_indicator,omitempty"`
	CreationDatetime               *ScotiaDatetime             `json:"creation_date_time,omitempty"`
	PaymentExpiryDate              *ScotiaDatetime             `json:"payment_expiry_date,omitempty"`
	SuppressResponderNotifications bool                        `json:"suppress_responder_notifications,omitempty"`
	ReturnUrl                      string                      `json:"return_url,omitempty"` //Need?
	Language                       string                      `json:"language,omitempty"`
	InstructedAmtData              *ScotiaAmtData              `json:"instructed_amount,omitempty"`
	InitiatingParty                *InitiatingPartyData        `json:"initiating_party,omitempty"`
	Debtor                         *DebtorData                 `json:"debtor,omitempty"`
	UltimateDebtor                 *DebtorData                 `json:"ultimate_debtor,omitempty"`
	Creditor                       *CreditorData               `json:"creditor,omitempty"`
	UltimateCreditor               *CreditorData               `json:"ultimate_creditor,omitempty"`
	CreditorAccount                *CreditorAccountData        `json:"creditor_account,omitempty"`
	FraudSupplementaryInfo         *FraudSupplementaryInfoData `json:"fraud_supplementary_info,omitempty"`
	PaymentCondition               *PaymentConditionData       `json:"payment_condition,omitempty"`
	PaymentTypeInformation         *PaymentTypeInformationData `json:"payment_type_information,omitempty"`
	RemittanceInformation          *RemittanceInformationData  `json:"remittance_information,omitempty"`
}

type RequestPaymentRequest struct {
	RequestData *RequestPaymentRequestData `json:"data,omitempty"`
}

type RequestPaymentResponseData struct {
	PaymentId               string `json:"payment_id,omitempty"`
	ClearingSystemReference string `json:"clearing_system_reference,omitempty"`
	Status                  string `json:"status,omitempty"` //not sure if this field exists
}

type RequestPaymentResponse struct {
	ResponseData *RequestPaymentResponseData `json:"data,omitempty"`
}
