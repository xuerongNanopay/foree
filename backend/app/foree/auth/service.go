package auth

import (
	"fmt"

	"xue.io/go-pay/auth"
)

type AuthService struct {
	sessionRepo            *auth.SessionRepo
	userRepo               *auth.UserRepo
	emailPasswordRepo      *auth.EmailPasswdRepo
	permissionRepo         *auth.PermissionRepo
	emailPasswdRecoverRepo *auth.EmailPasswdRecoverRepo
}

// Any error should return 503
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

	session = a.sessionRepo.GetSessionUniqueById(sessionId)
	if session == nil {
		return nil, fmt.Errorf("sesson `%s` not found", sessionId)
	}
	return session, nil
}

func (a *AuthService) VerifyEmail(session *auth.Session, req VerifyEmailReq) (*auth.Session, error) {
	return nil, nil
}

func (a *AuthService) Login() (*auth.Session, error) {
	return nil, nil
}

func (a *AuthService) CreateUser() (*auth.Session, error) {
	return nil, nil
}

func (a *AuthService) ResendVerifyCode() {

}

func (a *AuthService) ForgetPassword(email string) {

}

func (a *AuthService) ForgetPasswordUpdate(code, newPassword string) {

}

func (a *AuthService) Logout(session *auth.Session) {
	a.sessionRepo.Delete(session.ID)
}

func (a *AuthService) GetUser() {

}

func (a *AuthService) GetSession(sessionId string) *auth.Session {
	return nil
}

func (a *AuthService) Authorize(session auth.Session, permission string) bool {
	return false
}
