package foree_router

import (
	"github.com/gorilla/mux"
	foree_account_service "xue.io/go-pay/app/foree/service/account"
	foree_auth_service "xue.io/go-pay/app/foree/service/auth"
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
			"Foree:AcountService:VerifyContact",
			c.authService,
			c.accountService.VerifyContact,
		),
	).Methods("POST")
	// router.HandleFunc("/create_contact_account", simplePostWrapper(c.accountService.CreateContact)).Methods("POST")
	router.HandleFunc(
		"/create_contact_account",
		sessionPostWrapper(
			"VerifyContact",
			"Foree:AcountService:CreateContact",
			c.authService,
			c.accountService.CreateContact,
		),
	).Methods("POST")
	// router.HandleFunc("/delete_contact_account", simplePostWrapper(c.accountService.DeleteContact)).Methods("POST")
	router.HandleFunc(
		"/delete_contact_account",
		sessionPostWrapper(
			"DeleteContact",
			"Foree:AcountService:DeleteContact",
			c.authService,
			c.accountService.DeleteContact,
		),
	).Methods("POST")
	// router.HandleFunc("/contact_accounts/{ContactId}", simpleGetWrapper(c.accountService.GetActiveContact)).Methods("GET")
	router.HandleFunc(
		"/contact_accounts/{ContactId}",
		sessionGetWrapper(
			"GetActiveContact",
			"Foree:AcountService:GetActiveContact",
			c.authService,
			c.accountService.GetActiveContact,
		),
	).Methods("GET")
	// router.HandleFunc("/contact_accounts", simpleGetWrapper(c.accountService.GetAllActiveContacts)).Methods("GET")
	router.HandleFunc(
		"/contact_accounts",
		sessionGetWrapper(
			"GetAllActiveContacts",
			"Foree:AcountService:GetAllActiveContacts",
			c.authService,
			c.accountService.GetAllActiveContacts,
		),
	).Methods("GET")
	// router.HandleFunc("/interac_accounts", simpleGetWrapper(c.accountService.GetAllActiveInteracs)).Methods("GET")
	router.HandleFunc(
		"/interac_accounts",
		sessionGetWrapper(
			"GetAllActiveInteracs",
			"Foree:AcountService:GetAllActiveInteracs",
			c.authService,
			c.accountService.GetAllActiveInteracs,
		),
	).Methods("GET")
}
