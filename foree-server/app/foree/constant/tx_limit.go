package foree_constant

type TransactionLimitGroup string

const (
	TLPersonal1k TransactionLimitGroup = "TRANSACTION_PERSONAL_LIMIT_1K"
	TLPersonal2k TransactionLimitGroup = "TRANSACTION_PERSONAL_LIMIT_2K"
	TLPersonal3k TransactionLimitGroup = "TRANSACTION_PERSONAL_LIMIT_3K"
	TLBO         TransactionLimitGroup = "TRANSACTION_BO_LIMIT"
	TLAdmin      TransactionLimitGroup = "TRANSACTION_ADMIN_LIMIT"
)
