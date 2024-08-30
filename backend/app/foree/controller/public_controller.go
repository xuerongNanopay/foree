package foree_controller

import (
	"github.com/gorilla/mux"
	"xue.io/go-pay/app/foree/service"
)

type PublicController struct {
	transactionService *service.TransactionService
}

func NewPublicController(transactionService *service.TransactionService) *PublicController {
	return &PublicController{
		transactionService: transactionService,
	}
}

func (c *PublicController) RegisterRouter(router *mux.Router) {

}
