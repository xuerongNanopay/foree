package foree_controller

import (
	"github.com/gorilla/mux"
	"xue.io/go-pay/app/foree/service"
)

type PublicController struct {
}

func NewPublicController(transactionService *service.TransactionService) *PublicController {
	return &PublicController{}
}

func (c *PublicController) RegisterRouter(router *mux.Router) {

}
