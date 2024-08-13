package service

import (
	"fmt"
	"strings"

	"xue.io/go-pay/app/foree/transport"
)

type FreeQuoteReq struct {
	SrcAmount    float64 `json:"srcAmount"`
	SrcCurrency  string  `json:"srcCurrency" validate:"required"`
	DestAmount   float64 `json:"DestAmount"`
	DestCurrency string  `json:"DestCurrency" validate:"required"`
}

func (q *FreeQuoteReq) TrimSpace() {
	q.SrcCurrency = strings.TrimSpace(q.SrcCurrency)
	q.DestCurrency = strings.TrimSpace(q.DestCurrency)
}

func (q *FreeQuoteReq) Validate() *transport.BadRequestError {
	q.TrimSpace()
	ret := validateStruct(q, "Invalid free quote request")

	if q.SrcAmount <= 0 && q.DestAmount <= 0 {
		ret.AddDetails("srcAmount", fmt.Sprintf("invalid srcAmount `%v`", q.SrcAmount))
		ret.AddDetails("DestAmount", fmt.Sprintf("invalid DestAmount `%v`", q.DestAmount))
	}

	if len(ret.Details) > 0 {
		return ret
	}
	return nil
}

type QuoteTransactionReq struct {
}

type ConfirmQuoteReq struct {
	QuoteId string `json:"quoteId" validate:"required"`
}

func (q *ConfirmQuoteReq) TrimSpace() {
	q.QuoteId = strings.TrimSpace(q.QuoteId)
}

type GetTransactionReq struct {
}

type QueryTransactionReq struct {
}

type CancelTransactionReq struct {
}

type GetRateReq struct {
}
