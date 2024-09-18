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
		sessionGetWrapper(
			"QuoteTx",
			foree_service.PermissionForeeTxWrite,
			c.authService,
			c.transactionService.QuoteTx,
		),
	).Methods("POST")
	// Transaction creation
	router.HandleFunc(
		"/create_transaction",
		sessionGetWrapper(
			"QuoteTx",
			foree_service.PermissionForeeTxWrite,
			c.authService,
			c.transactionService.CreateTx,
		),
	).Methods("POST")
	// Transaction limit
	router.HandleFunc(
		"/transaction_limit",
		sessionGetWrapper(
			"QuoteTx",
			foree_service.PermissionForeeTxWrite,
			c.authService,
			c.transactionService.GetDailyTxLimit,
		),
	).Methods("GET")
	// Summary Transaction detail
	router.HandleFunc("/transactions/{TransactionId}", simpleGetWrapper(c.transactionService.GetTxSummary)).Methods("GET")
	// Summary Transaction query
	router.HandleFunc("/transactions", simpleGetWrapper(c.transactionService.QuerySummaryTxs)).Methods("GET")

}
