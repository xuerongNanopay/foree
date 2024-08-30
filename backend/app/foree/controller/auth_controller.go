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
	// Verify email
	router.HandleFunc("/verify_email", simplePostWrapper(c.authService.VerifyEmail)).Methods("POST")
	// Resend verify code
	router.HandleFunc("/resend_code", simpleGetWrapper(c.authService.VerifyEmail)).Methods("GET")
	// Signup
	router.HandleFunc("/sign_up", simplePostWrapper(c.authService.SignUp)).Methods("POST")
	// Logout
	router.HandleFunc("/logout", simpleGetWrapper(c.authService.Logout)).Methods("GET")
	// Onboard
	router.HandleFunc("/onboard", simplePostWrapper(c.authService.CreateUser)).Methods("POST")
}
