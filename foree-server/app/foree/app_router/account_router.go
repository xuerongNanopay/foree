package foree_router

import (
	"github.com/gorilla/mux"
	foree_service "xue.io/go-pay/app/foree/app_service"
)

type AccountRouter struct {
	accountService *foree_service.AccountService
}

func NewAccountRouter(accountService *foree_service.AccountService) *AccountRouter {
	return &AccountRouter{
		accountService: accountService,
	}
}

func (c *AccountRouter) RegisterRouter(router *mux.Router) {
	//TODO: handle to verify contact.
	router.HandleFunc("/verify_contact_account", simplePostWrapper(c.accountService.VerifyContact)).Methods("POST")
	router.HandleFunc("/create_contact_account", simplePostWrapper(c.accountService.CreateContact)).Methods("POST")
	router.HandleFunc("/delete_contact_account", simplePostWrapper(c.accountService.DeleteContact)).Methods("POST")
	router.HandleFunc("/contact_accounts/{ContactId}", simpleGetWrapper(c.accountService.GetActiveContact)).Methods("GET")
	router.HandleFunc("/contact_accounts", simpleGetWrapper(c.accountService.GetAllActiveContacts)).Methods("GET")
	router.HandleFunc("/interac_accounts", simpleGetWrapper(c.accountService.GetAllActiveInteracs)).Methods("GET")
}
