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
	router.HandleFunc("/forget_passwd", simplePostWrapper(c.authService.ForgetPasswd)).Methods("POST")
	// Forget password verify
	router.HandleFunc("/forget_passwd_verify", simplePostWrapper(c.authService.ForgetPasswdVerify)).Methods("POST")
	// Forget password update
	router.HandleFunc("/forget_passwd_update", simplePostWrapper(c.authService.ForgetPasswdUpdate)).Methods("POST")
	// Get user
	router.HandleFunc("/user", simpleGetWrapper(c.authService.GetUser)).Methods("GET")
	// Get user detail
	router.HandleFunc("/user_detail", simpleGetWrapper(c.authService.GetUserDetail)).Methods("GET")
	// Get user setting
	router.HandleFunc("/user_setting", simpleGetWrapper(c.authService.GetUserSetting)).Methods("GET")
	// Update passwd
	router.HandleFunc("/update_passwd", simplePostWrapper(c.authService.UpdatePasswd)).Methods("POST")
	// Update address
	router.HandleFunc("/update_address", simplePostWrapper(c.authService.UpdateUserAddress)).Methods("POST")
	// Update phone
	router.HandleFunc("/update_phone", simplePostWrapper(c.authService.UpdateUserPhoneNumber)).Methods("POST")
	// Update user setting
	router.HandleFunc("/update_user_setting", simplePostWrapper(c.authService.UpdateUserSetting)).Methods("POST")
}
