package foree_controller

import (
	"github.com/gorilla/mux"
	"xue.io/go-pay/app/foree/service"
)

type AuthController struct {
	authService *service.AuthService
}

func (c *AuthController) RegisterRouter(router *mux.Router) {
	// Login
	router.HandleFunc("/login", simplePostWrapper(c.authService.Login)).Methods("POST")
}
