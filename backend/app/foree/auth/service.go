package auth

import (
	"context"
	"time"

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

func (a *AuthService) SignUp(ctx context.Context, req SignUpReq) (*auth.Session, transport.ForeeError) {
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

	//TODO: send email. by goroutine

	session := a.sessionRepo.GetSessionUniqueById(sessionId)
	if session == nil {
		return nil, transport.NewInteralServerError("sesson `%s` not found", sessionId)
	}
	return session, nil
}

func (a *AuthService) allowVerifyEmail(sessionId string) (*auth.Session, transport.ForeeError) {
	session := a.sessionRepo.GetSessionUniqueById(sessionId)
	if session == nil || session.EmailPasswd == nil {
		return nil, transport.NewPreconditionRequireError(
			transport.PreconditionRequireMsgLogin,
			transport.RequireActionLogin,
		)
	}
	if session.EmailPasswd.Status != auth.EPStatusWaitingVerify {
		return nil, transport.NewPreconditionRequireError(
			transport.PreconditionRequireMsgToMain,
			transport.RequireActionToMain,
		)
	}

	return session, nil
}

func (a *AuthService) VerifyEmail(ctx context.Context, req VerifyEmailReq) (*auth.Session, transport.ForeeError) {
	// Check Allow to VerifyEmail
	session, err := a.allowVerifyEmail(req.SessionId)
	if err != nil {
		return nil, err
	}

	if session.EmailPasswd.VerifyCode != req.Code {
		return nil, transport.NewFormError("Invalid VerifyEmail Requst", "verify code", "Not equal")
	}

	// VerifyEmail and update EmailPasswd.
	newEP := *session.EmailPasswd
	newEP.Status = auth.EPStatusActive
	newEP.CodeVerifiedAt = time.Now()
	e := a.emailPasswordRepo.UpdateEmailPasswdByEmail(newEP)
	if e != nil {
		return nil, transport.WrapInteralServerError(e)
	}

	ep, e := a.emailPasswordRepo.GetUniqueEmailPasswdById(newEP.ID)
	if e != nil {
		return nil, transport.WrapInteralServerError(e)
	}
	if ep == nil {
		return nil, transport.NewInteralServerError("unable to get EmailPasswd with id: `%v`", newEP.ID)
	}

	// Update session
	newSession := *session
	newSession.EmailPasswd = ep

	session, e = a.sessionRepo.UpdateSession(newSession)
	if e != nil {
		return nil, transport.WrapInteralServerError(e)
	}
	return session, nil
}

func (a *AuthService) CreateUser(ctx context.Context) (*auth.Session, transport.ForeeError) {
	return nil, nil
}

func (a *AuthService) Login(ctx context.Context, req LoginReq) (*auth.Session, transport.ForeeError) {
	return nil, nil
}

func (a *AuthService) ResendVerifyCode(ctx context.Context, session SessionReq) {

}

func (a *AuthService) ForgetPassword(ctx context.Context, email string) {

}

func (a *AuthService) ForgetPasswordUpdate(ctx context.Context, req ForgetPasswordUpdateReq) {

}

func (a *AuthService) Logout(ctx context.Context, session SessionReq) transport.ForeeError {
	a.sessionRepo.Delete(session.SessionId)
	return transport.NewPreconditionRequireError(
		transport.PreconditionRequireMsgLogin,
		transport.RequireActionLogin,
	)
}

func (a *AuthService) GetUser(ctx context.Context, session SessionReq) {

}

func (a *AuthService) Authorize(ctx context.Context, sessionId string, permission string) (*auth.Session, transport.ForeeError) {
	session := a.sessionRepo.GetSessionUniqueById(sessionId)
	err := verifySession(session)
	if err != nil {
		return nil, err
	}
	for _, p := range session.Permissions {
		ok := auth.IsPermissionGrand(permission, p.Name)
		if ok {
			return session, nil
		}
	}
	return nil, transport.NewForbiddenError(permission)
}

func (a *AuthService) VerifySession(sctx context.Context, sessionId string) (*auth.Session, transport.ForeeError) {
	session := a.sessionRepo.GetSessionUniqueById(sessionId)
	err := verifySession(session)
	if err != nil {
		return nil, err
	}
	return session, nil
}
