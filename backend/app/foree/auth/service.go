package auth

import (
	"fmt"

	"xue.io/go-pay/app/foree/transport"
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
	// Check if email already exists.
	oldEmail, err := a.emailPasswordRepo.GetUniqueEmailPasswdByEmail(req.Email)
	if err != nil {
		return nil, err
	}

	if oldEmail != nil {
		return nil, transport.NewFormError(transport.FormErrorSignUpMsg, "email", "Duplicate email")
	}

	// Hashing password.
	hashedPassowrd, err := auth.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	// Create User
	userId, err := a.userRepo.InsertUser(auth.User{
		Status: auth.UserStatusInitial,
		Email:  req.Email,
		Group:  UserGroup,
	})

	if err != nil {
		return nil, err
	}

	user, err := a.userRepo.GetUniqueUserById(userId)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, fmt.Errorf("unable to get user with id: `%v`", userId)
	}

	// Create EmailPasswd
	id, err := a.emailPasswordRepo.InsertEmailPasswd(auth.EmailPasswd{
		Email:      req.Email,
		Passowrd:   hashedPassowrd,
		Status:     auth.EPStatusWaitingVerify,
		VerifyCode: auth.GenerateVerifyCode(),
		UserId:     user.ID,
	})

	if err != nil {
		return nil, err
	}

	ep, err := a.emailPasswordRepo.GetUniqueEmailPasswdById(id)

	if err != nil {
		return nil, err
	}

	if ep == nil {
		return nil, fmt.Errorf("unable to get EmailPasswd with id: `%v`", id)
	}

	sessionId, err := a.sessionRepo.InsertSession(&auth.Session{
		UserId:      user.ID,
		EmailPasswd: ep,
	})
	if err != nil {
		return nil, err
	}

	//TODO: send email.

	session := a.sessionRepo.GetSessionUniqueById(sessionId)
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
