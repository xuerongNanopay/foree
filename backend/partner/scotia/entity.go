package scotia

type RequestPayment struct {
	Paymentdata *RequestPaymentData `json:"data,omitempty"`
}

type RequestPaymentData struct {
}
