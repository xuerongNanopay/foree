package service

import (
	"fmt"
	"strings"

	"xue.io/go-pay/app/foree/transport"
)

type FreeQuoteReq struct {
	SrcAmount    float64 `json:"srcAmount"`
	SrcCurrency  string  `json:"srcCurrency" validate:"eq=CAD"`
	DestAmount   float64 `json:"DestAmount"`
	DestCurrency string  `json:"DestCurrency" validate:"eq=PKR"`
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
	transport.SessionReq
	SrcAccId     int64   `json:"srcAccId" validate:"gt=0"`
	DestAccId    int64   `json:"destAccId" validate:"gt=0"`
	SrcAmount    float64 `json:"srcAmount"`
	SrcCurrency  string  `json:"srcCurrency" validate:"eq=CAD"`
	DestAmount   float64 `json:"DestAmount"`
	DestCurrency string  `json:"DestCurrency" validate:"eq=PKR"`
}

func (q *QuoteTransactionReq) TrimSpace() {
	q.SrcCurrency = strings.TrimSpace(q.SrcCurrency)
	q.DestCurrency = strings.TrimSpace(q.DestCurrency)
}

func (q *QuoteTransactionReq) Validate() *transport.BadRequestError {
	q.TrimSpace()
	ret := validateStruct(q, "Invalid quote transaction request")

	if q.SrcAmount <= 0 && q.DestAmount <= 0 {
		ret.AddDetails("srcAmount", fmt.Sprintf("invalid srcAmount `%v`", q.SrcAmount))
		ret.AddDetails("DestAmount", fmt.Sprintf("invalid DestAmount `%v`", q.DestAmount))
	}

	if len(ret.Details) > 0 {
		return ret
	}
	return nil
}

type ConfirmQuoteReq struct {
	transport.SessionReq
	QuoteId string `json:"quoteId" validate:"required"`
}

func (q *ConfirmQuoteReq) TrimSpace() {
	q.QuoteId = strings.TrimSpace(q.QuoteId)
}

func (q *ConfirmQuoteReq) Validate() *transport.BadRequestError {
	q.TrimSpace()
	if ret := validateStruct(q, "Invalid confirm quote request"); len(ret.Details) > 0 {
		return ret
	}
	return nil
}

type GetTransactionReq struct {
	transport.SessionReq
	TransactionId int64 `json:"transactionId" validate:"required,gt=0"`
}

func (q *GetTransactionReq) TrimSpace() {
}

func (q *GetTransactionReq) Validate() *transport.BadRequestError {
	q.TrimSpace()
	if ret := validateStruct(q, "Invalid get transaction request"); len(ret.Details) > 0 {
		return ret
	}
	return nil
}

type QueryTransactionReq struct {
	transport.SessionReq
	Status string `json:""`
	Offset int    `json:"offset" validate:"required,gte=0"`
	Limit  int    `json:"limit" validate:"required,gt=0"`
}

func (q *QueryTransactionReq) TrimSpace() {
}

func (q *QueryTransactionReq) Validate() *transport.BadRequestError {
	q.TrimSpace()

	ret := validateStruct(q, "Invalid query transaction request")

	// Check status
	_, ok := allowTransactionsStatus[q.Status]
	if !ok {
		ret.AddDetails("status", fmt.Sprintf("invalid status `%v`", q.Status))
	}

	if len(ret.Details) > 0 {
		return ret
	}
	return nil
}

type CancelTransactionReq struct {
	transport.SessionReq
	TransactionId int64 `json:"transactionId" validate:"required,gt=0"`
}

func (q *CancelTransactionReq) TrimSpace() {
}

func (q *CancelTransactionReq) Validate() *transport.BadRequestError {
	q.TrimSpace()
	if ret := validateStruct(q, "Invalid cancel transaction request"); len(ret.Details) > 0 {
		return ret
	}
	return nil
}

type GetRateReq struct {
	SrcCurrency  string `json:"srcCurrency" validate:"eq=CAD"`
	DestCurrency string `json:"DestCurrency" validate:"eq=PKR"`
}

func (q *GetRateReq) TrimSpace() {
	q.SrcCurrency = strings.TrimSpace(q.SrcCurrency)
	q.DestCurrency = strings.TrimSpace(q.DestCurrency)
}

func (q *GetRateReq) Validate() *transport.BadRequestError {
	q.TrimSpace()
	ret := validateStruct(q, "Invalid free quote request")

	if len(ret.Details) > 0 {
		return ret
	}
	return nil
}

// ----------   Response --------------
