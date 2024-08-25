package foree_constant

import (
	"xue.io/go-pay/app/foree/transaction"
	"xue.io/go-pay/app/foree/types"
)

type RoleGroup string

const (
	RoleGroupPersonal RoleGroup = "FOREE_PERSONAL"
	RoleGroupBO       RoleGroup = "FOREE_BO"
	RoleGroupAdmin    RoleGroup = "FOREE_ADMIN"
)

type TransactionLimitGroup string

const (
	TLPersonal1k TransactionLimitGroup = "TRANSACTION_PERSONAL_LIMIT_1K"
	TLPersonal2k TransactionLimitGroup = "TRANSACTION_PERSONAL_LIMIT_2K"
	TLPersonal3k TransactionLimitGroup = "TRANSACTION_PERSONAL_LIMIT_3K"
	TLBO         TransactionLimitGroup = "TRANSACTION_BO_LIMIT"
	TLAdmin      TransactionLimitGroup = "TRANSACTION_ADMIN_LIMIT"
)

const (
	FeeName           string = "FOREE_TX_CAD_FEE"
	DefaultForeeGroup string = "FOREE_PERSONAL"
)

// Group level transaction limit.
var TxLimits = map[TransactionLimitGroup]transaction.TxLimit{
	TLPersonal1k: {
		Name: string(TLPersonal1k),
		MinAmt: types.AmountData{
			Amount:   types.Amount(10.0),
			Currency: "CAD",
		},
		MaxAmt: types.AmountData{
			Amount:   types.Amount(1000.0),
			Currency: "CAD",
		},
		IsEnable: true,
	},
	TLPersonal2k: {
		Name: string(TLPersonal2k),
		MinAmt: types.AmountData{
			Amount:   types.Amount(10.0),
			Currency: "CAD",
		},
		MaxAmt: types.AmountData{
			Amount:   types.Amount(2000.0),
			Currency: "CAD",
		},
		IsEnable: true,
	},
	TLPersonal3k: {
		Name: string(TLPersonal2k),
		MinAmt: types.AmountData{
			Amount:   types.Amount(10.0),
			Currency: "CAD",
		},
		MaxAmt: types.AmountData{
			Amount:   types.Amount(3000.0),
			Currency: "CAD",
		},
		IsEnable: true,
	},
	TLBO: {
		Name: string(TLBO),
		MinAmt: types.AmountData{
			Amount:   types.Amount(1.0),
			Currency: "CAD",
		},
		MaxAmt: types.AmountData{
			Amount:   types.Amount(3000.0),
			Currency: "CAD",
		},
		IsEnable: true,
	},
}
