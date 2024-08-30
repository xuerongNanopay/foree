package controller

import (
	"fmt"
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
	// Login
	loginHandler := restful_wrapper.RestPostWrapper(
		c.authService.Login,
		func(w http.ResponseWriter, session *auth.Session) http.ResponseWriter {
			//TODO: add session
			return w
		},
		func(req service.LoginReq, session *auth.Session, hErr transport.HError) {
			if v, is := hErr.(*transport.InteralServerError); is {
				// use logger.
				fmt.Print(v.OriginalError.Error())
			} else {
				fmt.Println(hErr.Error())
			}
		},
		true,
	)

	router.HandleFunc("/login", loginHandler).Methods("POST")
}
