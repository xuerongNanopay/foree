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
	userIdentificationRepo *UserIdentificationRepo
	// emailPasswdRecoverRepo *auth.EmailPasswdRecoverRepo
}

// Any error should return 503

func (a *AuthService) SignUp(ctx context.Context, req SignUpReq) (*auth.Session, transport.ForeeError) {
	// Check if email already exists.
	oldEmail, err := a.emailPasswordRepo.GetUniqueEmailPasswdByEmail(req.Email)
	if err != nil {
		return nil, transport.WrapInteralServerError(err)
	}

	if oldEmail != nil {
		return nil, transport.NewFormError("Invaild Signup Request", "email", "Duplicate email")
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
		User:        user,
		EmailPasswd: ep,
	})
	if err != nil {
		return nil, transport.WrapInteralServerError(err)
	}

	//TODO: Update referral
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
		return nil, transport.NewFormError("Invalid VerifyEmail Requst", "verify code", "Do not match")
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

func (a *AuthService) ResendVerifyCode(ctx context.Context, req transport.SessionReq) transport.ForeeError {
	// Check Allow to VerifyEmail
	session, err := a.allowVerifyEmail(req.SessionId)
	if err != nil {
		return err
	}

	// Change VerifyCode
	newEP := *session.EmailPasswd
	newEP.VerifyCode = auth.GenerateVerifyCode()

	e := a.emailPasswordRepo.UpdateEmailPasswdByEmail(newEP)
	if e != nil {
		return transport.WrapInteralServerError(e)
	}

	ep, e := a.emailPasswordRepo.GetUniqueEmailPasswdById(newEP.ID)
	if e != nil {
		return transport.WrapInteralServerError(e)
	}
	if ep == nil {
		return transport.NewInteralServerError("unable to get EmailPasswd with id: `%v`", newEP.ID)
	}

	//TODO: send email. by goroutine

	// Update session
	newSession := *session
	newSession.EmailPasswd = ep

	_, e = a.sessionRepo.UpdateSession(newSession)
	if e != nil {
		return transport.WrapInteralServerError(e)
	}
	return nil
}

func (a *AuthService) allowCreateUser(sessionId string) (*auth.Session, transport.ForeeError) {
	session := a.sessionRepo.GetSessionUniqueById(sessionId)
	if session == nil || session.EmailPasswd == nil {
		return nil, transport.NewPreconditionRequireError(
			transport.PreconditionRequireMsgLogin,
			transport.RequireActionLogin,
		)
	}

	if session.EmailPasswd.Status == auth.EPStatusWaitingVerify {
		return nil, transport.NewPreconditionRequireError(
			transport.PreconditionRequireMsgVerifyEmail,
			transport.RequireActionVerifyEmail,
		)
	}

	if session.User != nil && session.User.Status == auth.UserStatusInitial {
		return nil, transport.NewPreconditionRequireError(
			transport.PreconditionRequireMsgToMain,
			transport.RequireActionToMain,
		)
	}

	return session, nil
}

// TODO: create default contact account.
func (a *AuthService) CreateUser(ctx context.Context, req CreateUserReq) (*auth.Session, transport.ForeeError) {
	// Check allow to create user
	session, err := a.allowCreateUser(req.SessionId)
	if err != nil {
		return nil, err
	}

	// Create identification(Store Identification first)
	identification := UserIdentification{
		Status:  IdentificationStatusApproved,
		Type:    IdentificationType(req.IdentificationType),
		Value:   req.IdentificationValue,
		OwnerId: session.User.ID,
	}

	_, ier := a.userIdentificationRepo.InsertUserIdentification(identification)
	if ier != nil {
		return nil, transport.WrapInteralServerError(ier)
	}

	// Create a new user by updating essential fields.
	newUser := *session.User
	newUser.FirstName = req.FirstName
	newUser.MiddleName = req.MiddleName
	newUser.LastName = req.LastName
	newUser.Age = req.Age
	newUser.Dob = req.Dob.Time
	newUser.Nationality = req.Nationality
	newUser.Address1 = req.Address1
	newUser.Address2 = req.Address2
	newUser.City = req.City
	newUser.Province = req.Province
	newUser.Country = req.Country
	newUser.PostalCode = req.PostalCode
	newUser.PhoneNumber = req.PhoneNumber

	updateErr := a.userRepo.UpdateUserById(newUser)

	if updateErr != nil {
		return nil, transport.WrapInteralServerError(updateErr)
	}

	user, er := a.userRepo.GetUniqueUserById(newUser.ID)
	if er != nil {
		return nil, transport.WrapInteralServerError(er)
	}

	// Get Permission.
	pers, pErr := a.permissionRepo.GetAllPermissionByGroupName(user.Group)
	if pErr != nil {
		return nil, transport.WrapInteralServerError(pErr)
	}

	// Update session.
	newSession := *session
	newSession.User = user
	newSession.Permissions = pers
	_, sessionErr := a.sessionRepo.UpdateSession(newSession)
	if sessionErr != nil {
		return nil, transport.WrapInteralServerError(sessionErr)
	}

	return &newSession, nil
}

func (a *AuthService) Login(ctx context.Context, req LoginReq) (*auth.Session, transport.ForeeError) {
	// Delete previous token if exists.
	a.sessionRepo.Delete(req.SessionId)

	// Verify email and password
	ep, err := a.emailPasswordRepo.GetUniqueEmailPasswdByEmail(req.Email)
	if err != nil {
		return nil, transport.WrapInteralServerError(err)
	}
	if ep == nil {
		return nil, transport.NewFormError("Invaild Signup Request", "email", "Invalid email")
	}

	ok := auth.ComparePasswords(req.Password, []byte(ep.Passowrd))
	if !ok {
		return nil, transport.NewFormError("Invaild Signup Request", "password", "Invalid password")
	}

	// Load user(user must exist, but not necessary to be active)
	user, err := a.userRepo.GetUniqueUserById(ep.UserId)
	if err != nil {
		return nil, transport.WrapInteralServerError(err)
	}
	if user == nil {
		return nil, transport.NewInteralServerError("User `%v` do not exists", ep.UserId)
	}

	// Load permissions
	pers, pErr := a.permissionRepo.GetAllPermissionByGroupName(user.Group)
	if pErr != nil {
		return nil, transport.WrapInteralServerError(pErr)
	}

	// Load Ip and User agent, and create session
	newSession := auth.Session{
		User:        user,
		UserId:      user.ID,
		EmailPasswd: ep,
		Permissions: pers,
	}

	ip, ok := ctx.Value("ip").(string)
	if ok {
		newSession.Ip = ip
	}

	userAgent, ok := ctx.Value("userAgent").(string)
	if ok {
		newSession.UserAgent = userAgent
	}

	sessionId, err := a.sessionRepo.InsertSession(newSession)
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

// func (a *AuthService) ForgetPassword(ctx context.Context, email string) {

// }

// func (a *AuthService) ForgetPasswordUpdate(ctx context.Context, req ForgetPasswordUpdateReq) {

// }

func (a *AuthService) Logout(ctx context.Context, session transport.SessionReq) transport.ForeeError {
	a.sessionRepo.Delete(session.SessionId)
	return transport.NewPreconditionRequireError(
		transport.PreconditionRequireMsgLogin,
		transport.RequireActionLogin,
	)
}

func (a *AuthService) GetUser(ctx context.Context, req transport.SessionReq) (*UserDTO, transport.ForeeError) {
	session, sErr := a.VerifySession(ctx, req.SessionId)
	if sErr != nil {
		return nil, sErr
	}

	return NewUserDTO(session.User), nil
}

func (a *AuthService) ChangePasswd(ctx context.Context, req ChangePasswdReq) transport.ForeeError {
	session, err := a.VerifySession(ctx, req.SessionId)
	if err != nil {
		return err
	}

	hashed, hErr := auth.HashPassword(req.Password)
	if hErr != nil {
		return transport.WrapInteralServerError(hErr)
	}
	ep := *session.EmailPasswd
	ep.Passowrd = hashed
	//TODO: log
	updateErr := a.emailPasswordRepo.UpdateEmailPasswdByEmail(ep)
	if updateErr != nil {
		return transport.WrapInteralServerError(updateErr)
	}
	return nil
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

func (a *AuthService) VerifySession(ctx context.Context, sessionId string) (*auth.Session, transport.ForeeError) {
	session := a.sessionRepo.GetSessionUniqueById(sessionId)
	err := verifySession(session)
	if err != nil {
		return nil, err
	}
	return session, nil
}
