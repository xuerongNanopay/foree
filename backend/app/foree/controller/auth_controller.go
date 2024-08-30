package foree_controller

import (
	"github.com/gorilla/mux"
	"xue.io/go-pay/app/foree/service"
	"xue.io/go-pay/auth"
	"xue.io/go-pay/server/restful_wrapper"
)

type AuthController struct {
	authService *service.AuthService
}

func (c *AuthController) RegisterRouter(router *mux.Router) {
	// Login
	loginHandler := restful_wrapper.RestPostWrapper(
		c.authService.Login,
		emptyBeforeResponse,
		afterLogger[*auth.Session],
		true,
	)

	router.HandleFunc("/login", loginHandler).Methods("POST")
}
