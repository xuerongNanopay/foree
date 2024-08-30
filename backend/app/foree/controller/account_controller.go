package foree_controller

import (
	"github.com/gorilla/mux"
	"xue.io/go-pay/app/foree/service"
)

type AccountController struct {
	accountService *service.AccountService
}

func NewAccountController(accountService *service.AccountService) *AccountController {
	return &AccountController{
		accountService: accountService,
	}
}

func (c *AccountController) RegisterRouter(router *mux.Router) {

}
