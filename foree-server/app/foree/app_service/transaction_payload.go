package foree_service

import (
	"fmt"

	"xue.io/go-pay/app/foree/account"
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
	SrcAmount          float64 `json:"srcAmount" validate:"gt=0"`
	SrcCurrency        string  `json:"srcCurrency" validate:"eq=CAD"`
	DestCurrency       string  `json:"destCurrency" validate:"eq=PKR"`
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
	TransactionId int64  `json:"transactionId" validate:"required,gt=0"`
	CancelReason  string `json:"cancelReason"`
}

func (q CancelTransactionReq) Validate() *transport.BadRequestError {
	if ret := validateStruct(q, "Invalid cancel transaction request"); len(ret.Details) > 0 {
		return ret
	}
	return nil
}

type GetRateReq struct {
	SrcCurrency  string `json:"srcCurrency" validate:"eq=CAD"`
	DestCurrency string `json:"destCurrency" validate:"eq=PKR"`
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
	SrcAmount    types.Amount `json:"srcAmount,omitempty"`
	SrcCurrency  string       `json:"srcCurrency,omitempty"`
	DestAmount   types.Amount `json:"destAmount,omitempty"`
	DestCurrency string       `json:"destCurrency,omitempty"`
	Description  string       `json:"description,omitempty"`
}

func NewRateDTO(r *transaction.Rate) *RateDTO {
	return &RateDTO{
		SrcAmount:    r.SrcAmt.Amount,
		SrcCurrency:  r.SrcAmt.Currency,
		DestAmount:   r.DestAmt.Amount,
		DestCurrency: r.DestAmt.Currency,
		Description:  r.ToSummary(),
	}
}

type DailyTxLimitDTO struct {
	UsedAmount   types.Amount `json:"usedAmount"`
	UsedCurrency string       `json:"usedCurrency,omitempty"`
	MaxAmount    types.Amount `json:"maxAmount"`
	MaxCurrency  string       `json:"maxCurrency,omitempty"`
}

func NewDailyTxLimitDTO(r *transaction.DailyTxLimit) *DailyTxLimitDTO {
	return &DailyTxLimitDTO{
		UsedAmount:   r.UsedAmt.Amount,
		UsedCurrency: r.UsedAmt.Currency,
		MaxAmount:    r.MaxAmt.Amount,
		MaxCurrency:  r.MaxAmt.Currency,
	}
}

type TxCancelDTO struct {
	TransactionId int64  `json:"transactionId,omitempty"`
	Message       string `json:"message,omitempty"`
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
	CreateAt        int64                       `json:"createAt,omitempty"`
}

func NewTxSummaryDTO(tx *transaction.TxSummary) *TxSummaryDTO {
	ret := &TxSummaryDTO{
		ID:              tx.ID,
		Summary:         tx.Summary,
		Type:            tx.Type,
		Status:          tx.Status,
		Rate:            tx.Rate,
		TotalAmount:     tx.Rate,
		TotalCurrency:   tx.TotalCurrency,
		IsCancelAllowed: tx.IsCancelAllowed,
	}

	if tx.CreatedAt == nil {
		ret.CreateAt = tx.CreatedAt.UnixMilli()
	}

	return ret
}

type TxSummaryDetailDTO struct {
	ID              int64                       `json:"id,omitempty"`
	Summary         string                      `json:"sumary,omitempty"`
	Type            string                      `json:"type,omitempty"`
	Status          transaction.TxSummaryStatus `json:"status,omitempty"`
	Rate            string                      `json:"rate,omitempty"`
	PaymentUrl      string                      `json:"paymentUrl,omitempty"`
	SrcAccSummary   string                      `json:"srcAccSummary"`
	SrcAmount       types.Amount                `json:"srcAmount"`
	SrcCurrency     string                      `json:"srcCurrency"`
	DestAccSummary  string                      `json:"destAccSummary"`
	DestAmount      types.Amount                `json:"destAmount"`
	DestCurrency    string                      `json:"destCurrency"`
	TotalAmount     types.Amount                `json:"totalAmount"`
	TotalCurrency   string                      `json:"totalCurrency"`
	FeeAmount       types.Amount                `json:"feeAmount"`
	FeeCurrency     string                      `json:"feeCurrency"`
	RewardAmount    types.Amount                `json:"rewardAmount"`
	RewardCurrency  string                      `json:"rewardCurrency"`
	IsCancelAllowed bool                        `json:"isCancelAllowed"`
	CreateAt        int64                       `json:"createAt,omitempty"`
	SrcAccount      *SumInteracAccountDTO       `json:"srcAccount,omitempty"`
	DestAccount     *SumContactAccountDTO       `json:"destAccount,omitempty"`
}

func NewSumInteracAccountDTO(acc *account.InteracAccount) *SumInteracAccountDTO {
	return &SumInteracAccountDTO{
		ID:          acc.ID,
		FirstName:   acc.FirstName,
		MiddleName:  acc.MiddleName,
		LastName:    acc.LastName,
		Address1:    acc.Address1,
		Address2:    acc.Address2,
		City:        acc.City,
		Province:    acc.Province,
		Country:     acc.Country,
		PostalCode:  acc.PostalCode,
		PhoneNumber: acc.PhoneNumber,
		Email:       acc.Email,
	}
}

type SumInteracAccountDTO struct {
	ID          int64  `json:"id,omitempty"`
	FirstName   string `json:"firstName,omitempty"`
	MiddleName  string `json:"middleName,omitempty"`
	LastName    string `json:"lastName,omitempty"`
	Address1    string `json:"address1,omitempty"`
	Address2    string `json:"address2,omitempty"`
	City        string `json:"city,omitempty"`
	Province    string `json:"province,omitempty"`
	Country     string `json:"country,omitempty"`
	PostalCode  string `json:"postalCode,omitempty"`
	PhoneNumber string `json:"phoneNumber,omitempty"`
	Email       string `json:"email,omitempty"`
}

func NewSumContactAccountDTO(acc *account.ContactAccount) *SumContactAccountDTO {
	return &SumContactAccountDTO{
		ID:              acc.ID,
		Type:            acc.Type,
		FirstName:       acc.FirstName,
		MiddleName:      acc.MiddleName,
		LastName:        acc.LastName,
		Address1:        acc.Address1,
		Address2:        acc.Address2,
		City:            acc.City,
		Province:        acc.Province,
		Country:         acc.Country,
		PostalCode:      acc.PostalCode,
		PhoneNumber:     acc.PhoneNumber,
		InstitutionName: acc.InstitutionName,
		BranchNumber:    acc.BranchNumber,
		AccountNumber:   acc.AccountNumber,
	}
}

type SumContactAccountDTO struct {
	ID              int64                      `json:"id"`
	Status          account.AccountStatus      `json:"status"`
	Type            account.ContactAccountType `json:"type"`
	FirstName       string                     `json:"firstName"`
	MiddleName      string                     `json:"middleName"`
	LastName        string                     `json:"lastName"`
	Address1        string                     `json:"address1"`
	Address2        string                     `json:"address2"`
	City            string                     `json:"city"`
	Province        string                     `json:"province"`
	Country         string                     `json:"country"`
	PostalCode      string                     `json:"postalCode"`
	PhoneNumber     string                     `json:"phoneNumber"`
	InstitutionName string                     `json:"institutionName"`
	BranchNumber    string                     `json:"branchNumber"`
	AccountNumber   string                     `json:"accountNumber"`
	AccountHash     string                     `json:"accountHash"`
}

func NewTxSummaryDetailDTO(tx *transaction.TxSummary) *TxSummaryDetailDTO {
	ret := &TxSummaryDetailDTO{
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
	}

	if tx.CreatedAt != nil {
		ret.CreateAt = tx.CreatedAt.UnixMilli()
	}

	if tx.SrcAccount != nil {
		ret.SrcAccount = NewSumInteracAccountDTO(tx.SrcAccount)
	}

	if tx.DestAccount != nil {
		ret.DestAccount = NewSumContactAccountDTO(tx.DestAccount)
	}

	return ret
}

type QuoteTransactionDTO struct {
	QuoteId string             `json:"quoteId,omitempty"`
	TxSum   TxSummaryDetailDTO `json:"txSum,omitempty"`
}
