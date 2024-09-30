package foree_router

import (
	"github.com/gorilla/mux"
	foree_service "xue.io/go-pay/app/foree/app_service"
)

type TransactionRouter struct {
	transactionService *foree_service.TransactionService
	authService        *foree_service.AuthService
}

func NewTransactionRouter(authService *foree_service.AuthService, transactionService *foree_service.TransactionService) *TransactionRouter {
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
			foree_service.PermissionForeeTxWrite,
			c.authService,
			c.transactionService.QuoteTx,
		),
	).Methods("POST")

	// Transaction creation
	router.HandleFunc(
		"/create_transaction",
		sessionPostWrapper(
			"CreateTx",
			foree_service.PermissionForeeTxWrite,
			c.authService,
			c.transactionService.CreateTx,
		),
	).Methods("POST")

	// Transaction limit
	router.HandleFunc(
		"/transaction_limit",
		sessionGetWrapper(
			"GetDailyTxLimit",
			foree_service.PermissionForeeTxWrite,
			c.authService,
			c.transactionService.GetDailyTxLimit,
		),
	).Methods("GET")

	// reward
	router.HandleFunc(
		"/transaction_reward",
		sessionGetWrapper(
			"GetReward",
			foree_service.PermissionForeeTxWrite,
			c.authService,
			c.transactionService.GetReward,
		),
	).Methods("GET")

	// Summary Transaction detail
	router.HandleFunc("/transactions/{TransactionId}", simpleGetWrapper(c.transactionService.GetTxSummary)).Methods("GET")
	router.HandleFunc(
		"/transactions/{TransactionId}",
		sessionGetWrapper(
			"GetTxSummary",
			foree_service.PermissionForeeTxSummaryRead,
			c.authService,
			c.transactionService.GetTxSummary,
		),
	).Methods("GET")
	// Summary Transaction query
	router.HandleFunc(
		"/transactions",
		sessionGetWrapper(
			"QuerySummaryTxs",
			foree_service.PermissionForeeTxSummaryRead,
			c.authService,
			c.transactionService.QuerySummaryTxs,
		),
	).Methods("GET")
	router.HandleFunc(
		"/transactions/length",
		sessionGetWrapper(
			"CountSummaryTxs",
			foree_service.PermissionForeeTxSummaryRead,
			c.authService,
			c.transactionService.CountSummaryTxs,
		),
	).Methods("GET")

}
