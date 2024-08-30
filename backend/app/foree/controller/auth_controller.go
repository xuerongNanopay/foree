package controller

import (
	"net/http"

	"github.com/gorilla/mux"
	"xue.io/go-pay/app/foree/service"
	"xue.io/go-pay/auth"
	"xue.io/go-pay/server/restful_wrapper"
	"xue.io/go-pay/server/transport"
)

type AuthController struct {
	authService *service.AuthService
}

func (c *AuthController) RegisterRouter(router *mux.Router) {
	loginHandler := restful_wrapper.RestPostWrapper(
		c.authService.Login,
		func(w http.ResponseWriter) http.ResponseWriter {
			return w
		},
		func(req service.LoginReq, session *auth.Session, hErr transport.HError) {},
		true,
	)

	router.HandleFunc("/login", loginHandler).Methods("POST")
}
