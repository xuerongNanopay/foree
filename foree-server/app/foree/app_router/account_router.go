package foree_router

import (
	"github.com/gorilla/mux"
	foree_account_service "xue.io/go-pay/app/foree/service/account"
	foree_auth_service "xue.io/go-pay/app/foree/service/auth"
	foree_tx_service "xue.io/go-pay/app/foree/service/transaction"
)

// Do I need additional struct this it?
type AccountRouter struct {
	accountService *foree_account_service.AccountService
	authService    *foree_auth_service.AuthService
}

func NewAccountRouter(authService *foree_auth_service.AuthService, accountService *foree_account_service.AccountService) *AccountRouter {
	return &AccountRouter{
		authService:    authService,
		accountService: accountService,
	}
}

func (c *AccountRouter) RegisterRouter(router *mux.Router) {
	//TODO: handle to verify contact.
	// router.HandleFunc("/verify_contact_account", simplePostWrapper(c.accountService.VerifyContact)).Methods("POST")
	router.HandleFunc(
		"/verify_contact_account",
		sessionPostWrapper(
			"VerifyContact",
			foree_tx_service.PermissionContactWrite,
			c.authService,
			c.accountService.VerifyContact,
		),
	).Methods("POST")
	// router.HandleFunc("/create_contact_account", simplePostWrapper(c.accountService.CreateContact)).Methods("POST")
	router.HandleFunc(
		"/create_contact_account",
		sessionPostWrapper(
			"VerifyContact",
			foree_tx_service.PermissionContactWrite,
			c.authService,
			c.accountService.CreateContact,
		),
	).Methods("POST")
	// router.HandleFunc("/delete_contact_account", simplePostWrapper(c.accountService.DeleteContact)).Methods("POST")
	router.HandleFunc(
		"/delete_contact_account",
		sessionPostWrapper(
			"DeleteContact",
			foree_tx_service.PermissionContactWrite,
			c.authService,
			c.accountService.DeleteContact,
		),
	).Methods("POST")
	// router.HandleFunc("/contact_accounts/{ContactId}", simpleGetWrapper(c.accountService.GetActiveContact)).Methods("GET")
	router.HandleFunc(
		"/contact_accounts/{ContactId}",
		sessionGetWrapper(
			"GetActiveContact",
			foree_tx_service.PermissionContactRead,
			c.authService,
			c.accountService.GetActiveContact,
		),
	).Methods("GET")
	// router.HandleFunc("/contact_accounts", simpleGetWrapper(c.accountService.GetAllActiveContacts)).Methods("GET")
	router.HandleFunc(
		"/contact_accounts",
		sessionGetWrapper(
			"GetAllActiveContacts",
			foree_tx_service.PermissionContactRead,
			c.authService,
			c.accountService.GetAllActiveContacts,
		),
	).Methods("GET")
	// router.HandleFunc("/interac_accounts", simpleGetWrapper(c.accountService.GetAllActiveInteracs)).Methods("GET")
	router.HandleFunc(
		"/interac_accounts",
		sessionGetWrapper(
			"GetAllActiveInteracs",
			foree_tx_service.PermissionContactRead,
			c.authService,
			c.accountService.GetAllActiveInteracs,
		),
	).Methods("GET")
}
