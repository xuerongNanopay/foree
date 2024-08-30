package foree_controller

import (
	"github.com/gorilla/mux"
	"xue.io/go-pay/app/foree/service"
)

type TransactionController struct {
	transactionService *service.TransactionService
}

func NewTransactionController(transactionService *service.TransactionService) *TransactionController {
	return &TransactionController{
		transactionService: transactionService,
	}
}

func (c *TransactionController) RegisterRouter(router *mux.Router) {
	// === Public
	// Rate
	router.HandleFunc("rate", simplePostWrapper(c.transactionService.GetRate)).Methods("POST")
	// Free Quote
	router.HandleFunc("free_quote", simplePostWrapper(c.transactionService.FreeQuote)).Methods("POST")

	// === Private
	// Transaction quote
	router.HandleFunc("quote", simplePostWrapper(c.transactionService.QuoteTx)).Methods("POST")
	// Transaction create_transaction
	router.HandleFunc("create_transaction", simplePostWrapper(c.transactionService.CreateTx)).Methods("POST")
}
