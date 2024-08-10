package auth

import "xue.io/go-pay/auth"

type AuthService struct {
	sessionRepo            *auth.SessionRepo
	userRepo               *auth.UserRepo
	emailPasswordRepo      *auth.EmailPasswdRepo
	permissionRepo         *auth.PermissionRepo
	emailPasswdRecoverRepo *auth.EmailPasswdRecoverRepo
}

func (a *AuthService) SignUp(req SignUpReq) (*auth.Session, error) {
	hashedPassowrd, err := auth.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}
	id, err := a.emailPasswordRepo.InsertEmailPasswd(auth.EmailPasswd{
		Email:      req.Email,
		Passowrd:   hashedPassowrd,
		Status:     auth.EPStatusWaitingVerify,
		VerifyCode: auth.GenerateVerifyCode(),
	})

	if err != nil {
		return nil, err
	}

	ep, err := a.emailPasswordRepo.GetUniqueEmailPasswdById(id)

	if err != nil {
		return nil, err
	}

	session := &auth.Session{
		EmailPasswd: ep,
	}

	sessionId, err := a.sessionRepo.InsertSession(session)
	if err != nil {
		return nil, err
	}

	//TODO: send email.
	return a.sessionRepo.GetSessionUniqueById(sessionId), nil
}

func (a *AuthService) ResendVerifyCode() {

}

func (a *AuthService) VerifyEmail(session *auth.Session) {

}

func (a *AuthService) ForgetPassword(email string) {

}

func (a *AuthService) ForgetPasswordUpdate(code, newPassword string) {

}

func (a *AuthService) Login() {

}

func (a *AuthService) Logout(sessionId string) {

}

func (a *AuthService) CreateUser() {

}

func (a *AuthService) GetUser() {

}

func (a *AuthService) GetSession(sessionId string) *auth.Session {
	return nil
}

func (a *AuthService) Authorize(session auth.Session, permission string) bool {
	return false
}
