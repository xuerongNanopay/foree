package foree_router

import (
	"github.com/gorilla/mux"
	"xue.io/go-pay/app/foree/service"
)

type AccountRouter struct {
	accountService *service.AccountService
}

func NewAccountRouter(accountService *service.AccountService) *AccountRouter {
	return &AccountRouter{
		accountService: accountService,
	}
}

func (c *AccountRouter) RegisterRouter(router *mux.Router) {
	router.HandleFunc("/create_contact_account", simplePostWrapper(c.accountService.CreateContact)).Methods("POST")
	router.HandleFunc("/contact_accounts/{ContactId}", simpleGetWrapper(c.accountService.GetActiveContact)).Methods("GET")
	router.HandleFunc("/contact_accounts", simpleGetWrapper(c.accountService.QueryActiveContacts)).Methods("GET")
	router.HandleFunc("/interac_accounts", simpleGetWrapper(c.accountService.GetAllActiveInteracs)).Methods("GET")
}
