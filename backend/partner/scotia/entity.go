package scotia

import "time"

type ScotiaDatetime time.Time

func (d ScotiaDatetime) MarshalJSON() ([]byte, error) {
	t := time.Time(d)
	s := t.Format(time.RFC3339)
	return []byte("\"" + s + "\""), nil
}

type RequestPayment struct {
	Paymentdata *RequestPaymentData `json:"data,omitempty"`
}

type RequestPaymentData struct {
	ProductCode            string `json:"product_code,omitempty"`
	MessageIdentification  string `json:"message_identification,omitempty"`
	EndToEndIdentification string `json:"end_to_end_identification,omitempty"`
	CreditDebitIndicator   string `json:"credit_debit_indicator,omitempty"`
}
