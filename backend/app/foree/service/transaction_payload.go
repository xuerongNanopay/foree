package service

import (
	"fmt"
	"strings"
	"time"

	"xue.io/go-pay/app/foree/transaction"
	"xue.io/go-pay/app/foree/transport"
	"xue.io/go-pay/app/foree/types"
)

type FreeQuoteReq struct {
	SrcAmount    float64 `json:"srcAmount"`
	SrcCurrency  string  `json:"srcCurrency" validate:"eq=CAD"`
	DestCurrency string  `json:"DestCurrency" validate:"eq=PKR"`
}

func (q *FreeQuoteReq) TrimSpace() {
	q.SrcCurrency = strings.TrimSpace(q.SrcCurrency)
	q.DestCurrency = strings.TrimSpace(q.DestCurrency)
}

func (q *FreeQuoteReq) Validate() *transport.BadRequestError {
	q.TrimSpace()
	ret := validateStruct(q, "Invalid free quote request")

	if q.SrcAmount <= 0 {
		ret.AddDetails("srcAmount", fmt.Sprintf("invalid srcAmount `%v`", q.SrcAmount))
	}

	if len(ret.Details) > 0 {
		return ret
	}
	return nil
}

type QuoteTransactionReq struct {
	transport.SessionReq
	CinAccId           int64   `json:"cinAccId" validate:"gt=0"`
	CoutAccId          int64   `json:"coutAccId" validate:"gt=0"`
	SrcAmount          float64 `json:"srcAmount" validate:"gt=10,lt=1000"`
	SrcCurrency        string  `json:"srcCurrency" validate:"eq=CAD"`
	DestCurrency       string  `json:"DestCurrency" validate:"eq=PKR"`
	RewardIds          []int64 `json:"rewardIds" validate:"max=1"`
	PromoCode          string  `json:"promoCode"`
	TransactionPurpose string  `json:"transactionPurpose" validate:"required"`
}

func (q *QuoteTransactionReq) TrimSpace() {
	q.SrcCurrency = strings.TrimSpace(q.SrcCurrency)
	q.DestCurrency = strings.TrimSpace(q.DestCurrency)
}

func (q *QuoteTransactionReq) Validate() *transport.BadRequestError {
	q.TrimSpace()
	ret := validateStruct(q, "Invalid quote transaction request")

	//TODO: support promoCode
	// if q.PromoCode != "" && len(q.RewardIds) > 0 {
	// 	ret.AddDetails("promoCode", "cannot apply promocode and reward together")
	// }

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

type GetAllTransactionReq struct {
	transport.SessionReq
	Offset int `json:"offset" validate:"required,gte=0"`
	Limit  int `json:"limit" validate:"required,gt=0"`
}

func (q *GetAllTransactionReq) TrimSpace() {
}

func (q *GetAllTransactionReq) Validate() *transport.BadRequestError {
	q.TrimSpace()

	ret := validateStruct(q, "Invalid get transactions request")

	if len(ret.Details) > 0 {
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
type RateDTO struct {
	SrcCurrency  string
	DestCurrency string
	amt          float64
}

func NewRateDTO(r *transaction.Rate) *RateDTO {
	return &RateDTO{
		SrcCurrency:  r.SrcAmt.Currency,
		DestCurrency: r.DestAmt.Currency,
		amt:          r.GetForwardRate(),
	}
}

type TxSummaryDTO struct {
	ID              int64     `json:"id"`
	Summary         string    `json:"sumary"`
	Type            string    `json:"type"`
	Status          string    `json:"status"`
	Rate            string    `json:"rate"`
	TotalAmount     string    `json:"totalAmount"`
	TotalCurrency   string    `json:"totalCurrency"`
	IsCancelAllowed bool      `json:"isCancelAllowed"`
	CreateAt        time.Time `json:"createAt"`
}

func NewTxSummaryDTO(tx *transaction.TxSummary) *TxSummaryDTO {
	return &TxSummaryDTO{
		ID:              tx.ID,
		Summary:         tx.Summary,
		Type:            tx.Type,
		Status:          tx.Status,
		Rate:            tx.Rate,
		TotalAmount:     tx.Rate,
		TotalCurrency:   tx.TotalCurrency,
		IsCancelAllowed: tx.IsCancelAllowed,
		CreateAt:        tx.CreateAt,
	}
}

type TxSummaryDetailDTO struct {
	ID              int64        `json:"id"`
	Summary         string       `json:"sumary"`
	Type            string       `json:"type"`
	Status          string       `json:"status"`
	Rate            string       `json:"rate"`
	SrcAccSummary   string       `json:"srcAccSummary"`
	SrcAmount       types.Amount `json:"srcAmount"`
	SrcCurrency     string       `json:"srcCurrency"`
	DestAccSummary  string       `json:"destAccSummary"`
	DestAmount      types.Amount `json:"destAmount"`
	DestCurrency    string       `json:"destCurrency"`
	TotalAmount     types.Amount `json:"totalAmount"`
	TotalCurrency   string       `json:"totalCurrency"`
	FeeAmount       types.Amount `json:"feeAmount"`
	FeeCurrency     string       `json:"feeCurrency"`
	RewardAmount    types.Amount `json:"rewardAmount"`
	RewardCurrency  string       `json:"rewardCurrency"`
	IsCancelAllowed bool         `json:"isCancelAllowed"`
	CreateAt        time.Time    `json:"createAt"`
}

func NewTxSummaryDetailDTO(tx *transaction.TxSummary) *TxSummaryDetailDTO {
	return &TxSummaryDetailDTO{
		ID:              tx.ID,
		Summary:         tx.Summary,
		Type:            tx.Type,
		Status:          tx.Status,
		Rate:            tx.Rate,
		SrcAccSummary:   tx.SrcAccSummary,
		SrcAmount:       tx.SrcAmount,
		SrcCurrency:     tx.SrcCurrency,
		DestAccSummary:  tx.DestAccSummary,
		DestAmount:      tx.DestAmount,
		DestCurrency:    tx.DestCurrency,
		TotalAmount:     tx.TotalAmount,
		TotalCurrency:   tx.TotalCurrency,
		FeeAmount:       tx.FeeAmount,
		FeeCurrency:     tx.FeeCurrency,
		RewardAmount:    tx.RewardAmount,
		RewardCurrency:  tx.RewardCurrency,
		IsCancelAllowed: tx.IsCancelAllowed,
		CreateAt:        tx.CreateAt,
	}
}
