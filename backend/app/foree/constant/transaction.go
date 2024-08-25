package foree_constant

import "xue.io/go-pay/app/foree/transaction"

var AllowTransactionsStatus = map[string]bool{
	string(transaction.TxSummaryStatusInitial):      true,
	string(transaction.TxSummaryStatusAwaitPayment): true,
	string(transaction.TxSummaryStatusInProgress):   true,
	string(transaction.TxSummaryStatusCompleted):    true,
	string(transaction.TxSummaryStatusCancelled):    true,
	string(transaction.TxSummaryStatusRefunding):    true,
	string(transaction.TxSummaryStatusRefunded):     true,
}
