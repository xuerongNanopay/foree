package foree_router

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
	router.HandleFunc("/create_contact_account", simplePostWrapper(c.accountService.CreateContact)).Methods("POST")
	router.HandleFunc("/contact_accounts/{ContactId}", simpleGetWrapper(c.accountService.GetActiveContact)).Methods("GET")
	router.HandleFunc("/contact_accounts", simpleGetWrapper(c.accountService.QueryActiveContacts)).Methods("GET")
	router.HandleFunc("/interac_accounts", simpleGetWrapper(c.accountService.GetAllActiveInteracs)).Methods("GET")
}
