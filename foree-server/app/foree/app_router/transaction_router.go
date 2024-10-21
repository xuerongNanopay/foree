package foree_router

import (
	"github.com/gorilla/mux"
	foree_auth_service "xue.io/go-pay/app/foree/service/auth"
	foree_tx_service "xue.io/go-pay/app/foree/service/transaction"
)

type TransactionRouter struct {
	transactionService *foree_tx_service.TransactionService
	authService        *foree_auth_service.AuthService
}

func NewTransactionRouter(authService *foree_auth_service.AuthService, transactionService *foree_tx_service.TransactionService) *TransactionRouter {
	return &TransactionRouter{
		authService:        authService,
		transactionService: transactionService,
	}
}

func (c *TransactionRouter) RegisterRouter(router *mux.Router) {
	// === Public
	// Rate
	router.HandleFunc("/rate", simplePostWrapper(c.transactionService.GetRate)).Methods("POST")
	// Free Quote
	router.HandleFunc("/free_quote", simplePostWrapper(c.transactionService.FreeQuote)).Methods("POST")

	// === Private
	// Transaction quote
	router.HandleFunc(
		"/quote",
		sessionPostWrapper(
			"QuoteTx",
			foree_tx_service.PermissionForeeTxWrite,
			c.authService,
			c.transactionService.QuoteTx,
		),
	).Methods("POST")

	// Transaction creation
	router.HandleFunc(
		"/create_transaction",
		sessionPostWrapper(
			"CreateTx",
			foree_tx_service.PermissionForeeTxWrite,
			c.authService,
			c.transactionService.CreateTx,
		),
	).Methods("POST")

	// Transaction Cancel
	router.HandleFunc(
		"/cancel_transaction",
		sessionPostWrapper(
			"CreateTx",
			foree_tx_service.PermissionForeeTxWrite,
			c.authService,
			c.transactionService.CancelTransaction,
		),
	).Methods("POST")

	// Transaction limit
	router.HandleFunc(
		"/transaction_limit",
		sessionGetWrapper(
			"GetDailyTxLimit",
			foree_tx_service.PermissionForeeTxWrite,
			c.authService,
			c.transactionService.GetDailyTxLimit,
		),
	).Methods("GET")

	// reward
	router.HandleFunc(
		"/transaction_reward",
		sessionGetWrapper(
			"GetReward",
			foree_tx_service.PermissionForeeTxWrite,
			c.authService,
			c.transactionService.GetReward,
		),
	).Methods("GET")

	// Summary Transaction detail
	router.HandleFunc(
		"/transaction/{TransactionId}",
		sessionGetWrapper(
			"GetTxSummary",
			foree_tx_service.PermissionForeeTxSummaryRead,
			c.authService,
			c.transactionService.GetTxSummary,
		),
	).Methods("GET")
	// Summary Transaction query
	router.HandleFunc(
		"/transactions",
		sessionGetWrapper(
			"QuerySummaryTxs",
			foree_tx_service.PermissionForeeTxSummaryRead,
			c.authService,
			c.transactionService.QuerySummaryTxs,
		),
	).Methods("GET")
	router.HandleFunc(
		"/transactions/length",
		sessionGetWrapper(
			"CountSummaryTxs",
			foree_tx_service.PermissionForeeTxSummaryRead,
			c.authService,
			c.transactionService.CountSummaryTxs,
		),
	).Methods("GET")

}
