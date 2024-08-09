package auth

import "xue.io/go-pay/auth"

type AuthService struct {
	sessionRepo            *auth.SessionRepo
	userRepo               *auth.UserRepo
	emailPasswordRepo      *auth.EmailPasswdRepo
	permissionRepo         *auth.PermissionRepo
	emailPasswdRecoverRepo *auth.EmailPasswdRecoverRepo
}

func (a *AuthService) SignUp(req SignUpReq) {
	a.userRepo.Insert(auth.User{})
}

func (a *AuthService) ResendVerifyCode() {

}

func (a *AuthService) VerifyEmail() {

}

func (a *AuthService) ForgetPassword(email string) {

}

func (a *AuthService) ForgetPasswordUpdate(code, newPassword string) {

}

func (a *AuthService) Login() {

}

func (a *AuthService) Logout() {

}

func (a *AuthService) CreateUser() {

}

func (a *AuthService) GetSession(sessionId string) *auth.Session {
	return nil
}

func (a *AuthService) Authorize(session auth.Session, permission string) bool {
	return false
}
