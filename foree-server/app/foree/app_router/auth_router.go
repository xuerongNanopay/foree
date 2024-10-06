package foree_router

import (
	"github.com/gorilla/mux"
	foree_service "xue.io/go-pay/app/foree/app_service"
)

type AuthRouter struct {
	authService *foree_service.AuthService
}

func NewAuthRouter(authService *foree_service.AuthService) *AuthRouter {
	return &AuthRouter{
		authService: authService,
	}
}

func (c *AuthRouter) RegisterRouter(router *mux.Router) {
	// Login
	router.HandleFunc("/login", simplePostWrapper(c.authService.Login)).Methods("POST")
	// Verify email
	router.HandleFunc("/verify_email", simplePostWrapper(c.authService.VerifyEmail)).Methods("POST")
	// Resend verify code
	router.HandleFunc("/resend_code", simpleGetWrapper(c.authService.ResendVerifyCode)).Methods("GET")
	// Signup
	router.HandleFunc("/sign_up", simplePostWrapper(c.authService.SignUp)).Methods("POST")
	// Logout
	router.HandleFunc("/logout", simpleGetWrapper(c.authService.Logout)).Methods("GET")
	// Onboard
	router.HandleFunc("/onboard", simplePostWrapper(c.authService.CreateUser)).Methods("POST")
	// Forget password
	router.HandleFunc("/forget_password", simplePostWrapper(c.authService.ForgetPassword)).Methods("POST")
	// Forget password verify
	router.HandleFunc("/forget_password_verify", simplePostWrapper(c.authService.ForgetPasswordVerify)).Methods("POST")
	// Forget password update
	router.HandleFunc("/forget_password_update", simplePostWrapper(c.authService.ForgetPasswordUpdate)).Methods("POST")
	// Get user
	router.HandleFunc("/user", simpleGetWrapper(c.authService.GetUser)).Methods("GET")
	// Update passwd
	router.HandleFunc("/change_password", simplePostWrapper(c.authService.ChangePasswd)).Methods("POST")
}
