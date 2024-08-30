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

}
