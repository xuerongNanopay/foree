package service

import (
	"fmt"
	"time"

	foree_constant "xue.io/go-pay/app/foree/constant"
	"xue.io/go-pay/app/foree/transaction"
	"xue.io/go-pay/app/foree/types"
	"xue.io/go-pay/server/transport"
)

type FreeQuoteReq struct {
	SrcAmount    float64 `json:"srcAmount"`
	SrcCurrency  string  `json:"srcCurrency" validate:"eq=CAD"`
	DestCurrency string  `json:"DestCurrency" validate:"eq=PKR"`
}

func (q FreeQuoteReq) Validate() *transport.BadRequestError {
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

func (q QuoteTransactionReq) Validate() *transport.BadRequestError {
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

type CreateTransactionReq struct {
	transport.SessionReq
	QuoteId string `json:"quoteId" validate:"required"`
}

func (q CreateTransactionReq) Validate() *transport.BadRequestError {
	if ret := validateStruct(q, "Invalid create quote request"); len(ret.Details) > 0 {
		return ret
	}
	return nil
}

type GetTransactionReq struct {
	transport.SessionReq
	TransactionId int64 `json:"transactionId" validate:"required,gt=0"`
}

func (q GetTransactionReq) Validate() *transport.BadRequestError {
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

func (q GetAllTransactionReq) Validate() *transport.BadRequestError {
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

func (q QueryTransactionReq) Validate() *transport.BadRequestError {

	ret := validateStruct(q, "Invalid query transaction request")

	// Check status
	_, ok := foree_constant.AllowTransactionsStatus[q.Status]
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

func (q CancelTransactionReq) Validate() *transport.BadRequestError {
	if ret := validateStruct(q, "Invalid cancel transaction request"); len(ret.Details) > 0 {
		return ret
	}
	return nil
}

type GetRateReq struct {
	SrcCurrency  string `json:"srcCurrency" validate:"eq=CAD"`
	DestCurrency string `json:"DestCurrency" validate:"eq=PKR"`
}

func (q GetRateReq) Validate() *transport.BadRequestError {
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
	ID              int64                       `json:"id,omitempty"`
	Summary         string                      `json:"sumary,omitempty"`
	Type            string                      `json:"type,omitempty"`
	Status          transaction.TxSummaryStatus `json:"status,omitempty"`
	Rate            string                      `json:"rate,omitempty"`
	TotalAmount     string                      `json:"totalAmount,omitempty"`
	TotalCurrency   string                      `json:"totalCurrency,omitempty"`
	IsCancelAllowed bool                        `json:"isCancelAllowed,omitempty"`
	CreatedAt       time.Time                   `json:"createdAt,omitempty"`
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
		CreatedAt:       tx.CreatedAt,
	}
}

type TxSummaryDetailDTO struct {
	ID              int64                       `json:"id,omitempty"`
	Summary         string                      `json:"sumary,omitempty"`
	Type            string                      `json:"type,omitempty"`
	Status          transaction.TxSummaryStatus `json:"status,omitempty"`
	Rate            string                      `json:"rate,omitempty"`
	PaymentUrl      string                      `json:"paymentUrl,omitempty"`
	SrcAccSummary   string                      `json:"srcAccSummary,omitempty"`
	SrcAmount       types.Amount                `json:"srcAmount,omitempty"`
	SrcCurrency     string                      `json:"srcCurrency,omitempty"`
	DestAccSummary  string                      `json:"destAccSummary,omitempty"`
	DestAmount      types.Amount                `json:"destAmount,omitempty"`
	DestCurrency    string                      `json:"destCurrency,omitempty"`
	TotalAmount     types.Amount                `json:"totalAmount,omitempty"`
	TotalCurrency   string                      `json:"totalCurrency,omitempty"`
	FeeAmount       types.Amount                `json:"feeAmount,omitempty"`
	FeeCurrency     string                      `json:"feeCurrency,omitempty"`
	RewardAmount    types.Amount                `json:"rewardAmount,omitempty"`
	RewardCurrency  string                      `json:"rewardCurrency,omitempty"`
	IsCancelAllowed bool                        `json:"isCancelAllowed,omitempty"`
	CreatedAt       time.Time                   `json:"createdAt,omitempty"`
}

func NewTxSummaryDetailDTO(tx transaction.TxSummary) *TxSummaryDetailDTO {
	return &TxSummaryDetailDTO{
		ID:              tx.ID,
		Summary:         tx.Summary,
		Type:            tx.Type,
		Status:          tx.Status,
		Rate:            tx.Rate,
		PaymentUrl:      tx.PaymentUrl,
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
		CreatedAt:       tx.CreatedAt,
	}
}

type QuoteTransactionDTO struct {
	QuoteId string             `json:"quoteId,omitempty"`
	TxSum   TxSummaryDetailDTO `json:"txSum,omitempty"`
}
