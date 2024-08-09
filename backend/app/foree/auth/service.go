package foree_auth

import "xue.io/go-pay/auth"

type AuthService struct {
	sessionRepo       *auth.SessionRepo
	userRepo          *auth.UserRepo
	emailPasswordRepo *auth.EmailPasswdRepo
	permissionRepo    *auth.PermissionRepo
}

//getSeesion
//Login
//SignUp
//Logout
//authorize
//Create User
