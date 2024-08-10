package auth

import (
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

func (a *AuthService) SignUp(req SignUpReq) (*auth.Session, transport.ForeeError) {
	// Check if email already exists.
	oldEmail, err := a.emailPasswordRepo.GetUniqueEmailPasswdByEmail(req.Email)
	if err != nil {
		return nil, transport.WrapInteralServerError(err)
	}

	if oldEmail != nil {
		return nil, transport.NewFormError(transport.FormErrorSignUpMsg, "email", "Duplicate email")
	}

	// Hashing password.
	hashedPassowrd, err := auth.HashPassword(req.Password)
	if err != nil {
		return nil, transport.WrapInteralServerError(err)
	}

	// Create User
	userId, err := a.userRepo.InsertUser(auth.User{
		Status: auth.UserStatusInitial,
		Email:  req.Email,
		Group:  UserGroup,
	})

	if err != nil {
		return nil, transport.WrapInteralServerError(err)
	}

	user, err := a.userRepo.GetUniqueUserById(userId)
	if err != nil {
		return nil, transport.WrapInteralServerError(err)
	}

	if user == nil {
		return nil, transport.NewInteralServerError("unable to get user with id: `%v`", userId)
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
		return nil, transport.WrapInteralServerError(err)
	}

	ep, err := a.emailPasswordRepo.GetUniqueEmailPasswdById(id)

	if err != nil {
		return nil, transport.WrapInteralServerError(err)
	}

	if ep == nil {
		return nil, transport.NewInteralServerError("unable to get EmailPasswd with id: `%v`", id)
	}

	sessionId, err := a.sessionRepo.InsertSession(auth.Session{
		UserId:      user.ID,
		EmailPasswd: ep,
	})
	if err != nil {
		return nil, transport.WrapInteralServerError(err)
	}

	//TODO: send email.

	session := a.sessionRepo.GetSessionUniqueById(sessionId)
	if session == nil {
		return nil, transport.NewInteralServerError("sesson `%s` not found", sessionId)
	}
	return session, nil
}

func (a *AuthService) VerifyEmail(session *auth.Session, req VerifyEmailReq) (*auth.Session, transport.ForeeError) {
	return nil, nil
}

func (a *AuthService) Login() (*auth.Session, transport.ForeeError) {
	return nil, nil
}

func (a *AuthService) CreateUser() (*auth.Session, transport.ForeeError) {
	return nil, nil
}

func (a *AuthService) ResendVerifyCode() {

}

func (a *AuthService) ForgetPassword(email string) {

}

func (a *AuthService) ForgetPasswordUpdate(code, newPassword string) {

}

func (a *AuthService) Logout(sessionId string) transport.ForeeError {
	a.sessionRepo.Delete(sessionId)
	return transport.NewPreconditionRequireError(
		transport.PreconditionRequireMsgLogin,
		transport.RequireActionLogin,
	)
}

func (a *AuthService) GetUser() {

}

func (a *AuthService) GetSession(sessionId string) *auth.Session {
	return nil
}

func (a *AuthService) Authorize(sessionId string, permission string) transport.ForeeError {
	session := a.sessionRepo.GetSessionUniqueById(sessionId)
	err := verifySession(session)
	if err != nil {
		return err
	}
	for _, p := range session.Permissions {
		ok := auth.IsPermissionGrand(permission, p.Name)
		if ok {
			return nil
		}
	}
	return transport.NewForbiddenError(permission)
}

func (a *AuthService) VerifySession(sessionId string) transport.ForeeError {
	session := a.sessionRepo.GetSessionUniqueById(sessionId)
	return verifySession(session)
}
