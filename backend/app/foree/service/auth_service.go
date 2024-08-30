package service

import (
	"context"
	"database/sql"
	"time"

	"xue.io/go-pay/app/foree/account"
	foree_auth "xue.io/go-pay/app/foree/auth"
	foree_constant "xue.io/go-pay/app/foree/constant"
	"xue.io/go-pay/auth"
	"xue.io/go-pay/server/transport"
)

type AuthService struct {
	db                     *sql.DB
	sessionRepo            *auth.SessionRepo
	userRepo               *auth.UserRepo
	emailPasswordRepo      *auth.EmailPasswdRepo
	rolePermissionRepo     *auth.RolePermissionRepo
	userIdentificationRepo *foree_auth.UserIdentificationRepo
	interacRepo            *account.InteracAccountRepo
	userGroupRepo          *auth.UserGroupRepo
	// emailPasswdRecoverRepo *auth.EmailPasswdRecoverRepo
}

// Any error should return 503
// TODO: DB Transaction.
func (a *AuthService) SignUp(ctx context.Context, req SignUpReq) (*auth.Session, transport.HError) {
	// Check if email already exists.
	oldEmail, err := a.emailPasswordRepo.GetUniqueEmailPasswdByEmail(req.Email)
	if err != nil {
		return nil, transport.WrapInteralServerError(err)
	}

	if oldEmail != nil {
		return nil, transport.NewFormError("Invaild signup", "email", "Duplicate email")
	}

	// Hashing password.
	hashedPasswd, err := auth.HashPassword(req.Password)
	if err != nil {
		return nil, transport.WrapInteralServerError(err)
	}

	dTx, err := a.db.Begin()
	if err != nil {
		dTx.Rollback()
		return nil, transport.WrapInteralServerError(err)
	}

	// Create User
	userId, err := a.userRepo.InsertUser(ctx, auth.User{
		Status: auth.UserStatusInitial,
		Email:  req.Email,
	})

	if err != nil {
		dTx.Rollback()
		return nil, transport.WrapInteralServerError(err)
	}

	user, err := a.userRepo.GetUniqueUserById(userId)
	if err != nil {
		dTx.Rollback()
		return nil, transport.WrapInteralServerError(err)
	}

	if user == nil {
		dTx.Rollback()
		return nil, transport.NewInteralServerError("unable to get user with id: `%v`", userId)
	}

	// Create EmailPasswd
	id, err := a.emailPasswordRepo.InsertEmailPasswd(ctx, auth.EmailPasswd{
		Email:      req.Email,
		Passwd:     hashedPasswd,
		Status:     auth.EPStatusWaitingVerify,
		VerifyCode: auth.GenerateVerifyCode(),
		OwnerId:    user.ID,
	})

	if err != nil {
		dTx.Rollback()
		return nil, transport.WrapInteralServerError(err)
	}

	ep, err := a.emailPasswordRepo.GetUniqueEmailPasswdById(id)

	if err != nil {
		dTx.Rollback()
		return nil, transport.WrapInteralServerError(err)
	}

	if ep == nil {
		dTx.Rollback()
		return nil, transport.NewInteralServerError("unable to get EmailPasswd with id: `%v`", id)
	}

	if err = dTx.Commit(); err != nil {
		return nil, transport.WrapInteralServerError(err)
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

func (a *AuthService) allowVerifyEmail(sessionId string) (*auth.Session, transport.HError) {
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

func (a *AuthService) VerifyEmail(ctx context.Context, req VerifyEmailReq) (*auth.Session, transport.HError) {
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
	// newEP.CodeVerifiedAt = time.Now()
	e := a.emailPasswordRepo.UpdateEmailPasswdByEmail(ctx, newEP)
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

func (a *AuthService) ResendVerifyCode(ctx context.Context, req transport.SessionReq) (*auth.Session, transport.HError) {
	// Check Allow to VerifyEmail
	session, err := a.allowVerifyEmail(req.SessionId)
	if err != nil {
		return nil, err
	}

	// Change VerifyCode
	newEP := *session.EmailPasswd
	newEP.VerifyCode = auth.GenerateVerifyCode()

	e := a.emailPasswordRepo.UpdateEmailPasswdByEmail(ctx, newEP)
	if e != nil {
		return nil, transport.WrapInteralServerError(e)
	}

	//TODO: send email. by goroutine

	// Update session
	newSession := *session
	newSession.EmailPasswd = &newEP

	_, e = a.sessionRepo.UpdateSession(newSession)
	if e != nil {
		return nil, transport.WrapInteralServerError(e)
	}
	return nil, nil
}

func (a *AuthService) allowCreateUser(sessionId string) (*auth.Session, transport.HError) {
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

func (a *AuthService) CreateUser(ctx context.Context, req CreateUserReq) (*auth.Session, transport.HError) {
	// Check allow to create user
	session, err := a.allowCreateUser(req.SessionId)
	if err != nil {
		return nil, err
	}

	dTx, dErr := a.db.Begin()
	if dErr != nil {
		dTx.Rollback()
		return nil, transport.WrapInteralServerError(dErr)
	}

	// Create identification(Store Identification first)
	identification := foree_auth.UserIdentification{
		Status:  foree_auth.IdentificationStatusActive,
		Type:    foree_auth.IdentificationType(req.IdentificationType),
		Value:   req.IdentificationValue,
		OwnerId: session.User.ID,
	}

	_, ier := a.userIdentificationRepo.InsertUserIdentification(ctx, identification)
	if ier != nil {
		dTx.Rollback()
		return nil, transport.WrapInteralServerError(ier)
	}

	// Create a new user by updating essential fields.
	newUser := *session.User
	newUser.FirstName = req.FirstName
	newUser.MiddleName = req.MiddleName
	newUser.LastName = req.LastName
	newUser.Age = req.Age
	newUser.Dob = req.Dob.Time
	// newUser.Nationality = req.Nationality
	newUser.Address1 = req.Address1
	newUser.Address2 = req.Address2
	newUser.City = req.City
	newUser.Province = req.Province
	newUser.Country = req.Country
	newUser.PostalCode = req.PostalCode
	newUser.PhoneNumber = req.PhoneNumber

	updateErr := a.userRepo.UpdateUserById(ctx, newUser)

	if updateErr != nil {
		dTx.Rollback()
		return nil, transport.WrapInteralServerError(updateErr)
	}

	user, er := a.userRepo.GetUniqueUserById(newUser.ID)
	if er != nil {
		dTx.Rollback()
		return nil, transport.WrapInteralServerError(er)
	}

	//Create userGroup
	_, er = a.userGroupRepo.InsertUserGroup(ctx, auth.UserGroup{
		RoleGroup:             foree_constant.DefaultRoleGroup,
		TransactionLimitGroup: foree_constant.DefaultTransactionLimitGroup,
		OwnerId:               user.ID,
	})

	userGroup, er := a.userGroupRepo.GetUniqueUserGroupByOwnerId(user.ID)
	if er != nil {
		dTx.Rollback()
		return nil, transport.WrapInteralServerError(er)
	}

	if err := dTx.Commit(); err != nil {
		return nil, transport.WrapInteralServerError(err)
	}

	// Get Permission.
	rolePermissions, pErr := a.rolePermissionRepo.GetAllEnabledRolePermissionByRoleName(userGroup.RoleGroup)
	if pErr != nil {
		return nil, transport.WrapInteralServerError(pErr)
	}

	// Update session.
	newSession := *session
	newSession.User = user
	newSession.UserId = user.ID
	newSession.RolePermissions = rolePermissions
	newSession.UserGroup = userGroup

	updateSession, sessionErr := a.sessionRepo.UpdateSession(newSession)
	if sessionErr != nil {
		return nil, transport.WrapInteralServerError(sessionErr)
	}

	// Create default Interac Account for the user.
	acc := account.InteracAccount{
		FirstName:        session.User.FirstName,
		MiddleName:       session.User.MiddleName,
		LastName:         session.User.LastName,
		Address1:         user.Address1,
		Address2:         user.Address2,
		City:             user.City,
		Province:         user.Province,
		Country:          user.Country,
		PostalCode:       user.PostalCode,
		PhoneNumber:      session.User.PhoneNumber,
		Email:            session.User.Email,
		OwnerId:          session.User.ID,
		Status:           account.AccountStatusActive,
		LatestActivityAt: time.Now(),
	}
	_, derr := a.interacRepo.InsertInteracAccount(ctx, acc)
	if derr != nil {
		return nil, transport.WrapInteralServerError(derr)
	}
	return updateSession, nil
}

func (a *AuthService) Login(ctx context.Context, req LoginReq) (*auth.Session, transport.HError) {
	// Delete previous token if exists.
	a.sessionRepo.Delete(req.SessionId)

	// Verify email and password
	ep, err := a.emailPasswordRepo.GetUniqueEmailPasswdByEmail(req.Email)
	if err != nil {
		return nil, transport.WrapInteralServerError(err)
	}
	if ep == nil {
		return nil, transport.NewFormError("Invaild signup", "email", "Invalid email")
	}

	ok := auth.ComparePasswords(req.Password, []byte(ep.Passwd))
	if !ok {
		return nil, transport.NewFormError("Invaild signup", "password", "Invalid password")
	}

	// Load user(user must exist, but not necessary to be active)
	user, err := a.userRepo.GetUniqueUserById(ep.OwnerId)
	if err != nil {
		return nil, transport.WrapInteralServerError(err)
	}
	if user == nil {
		return nil, transport.NewInteralServerError("User `%v` do not exists", ep.OwnerId)
	}

	userGroup, er := a.userGroupRepo.GetUniqueUserGroupByOwnerId(user.ID)
	if er != nil {
		return nil, transport.WrapInteralServerError(er)
	}

	// Load permissions
	pers, pErr := a.rolePermissionRepo.GetAllEnabledRolePermissionByRoleName(userGroup.RoleGroup)
	if pErr != nil {
		return nil, transport.WrapInteralServerError(pErr)
	}

	// Load Ip and User agent, and create session
	newSession := auth.Session{
		User:            user,
		UserId:          user.ID,
		UserGroup:       userGroup,
		EmailPasswd:     ep,
		RolePermissions: pers,
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

func (a *AuthService) Logout(ctx context.Context, session transport.SessionReq) (*auth.Session, transport.HError) {
	a.sessionRepo.Delete(session.SessionId)
	return nil, transport.NewPreconditionRequireError(
		transport.PreconditionRequireMsgLogin,
		transport.RequireActionLogin,
	)
}

func (a *AuthService) GetUser(ctx context.Context, req transport.SessionReq) (*UserDTO, transport.HError) {
	session, sErr := a.VerifySession(ctx, req.SessionId)
	if sErr != nil {
		return nil, sErr
	}

	return NewUserDTO(session.User), nil
}

func (a *AuthService) ChangePasswd(ctx context.Context, req ChangePasswdReq) transport.HError {
	session, err := a.VerifySession(ctx, req.SessionId)
	if err != nil {
		return err
	}

	hashed, hErr := auth.HashPassword(req.Password)
	if hErr != nil {
		return transport.WrapInteralServerError(hErr)
	}
	ep := *session.EmailPasswd
	ep.Passwd = hashed
	//TODO: log

	updateErr := a.emailPasswordRepo.UpdateEmailPasswdByEmail(ctx, ep)
	if updateErr != nil {
		return transport.WrapInteralServerError(updateErr)
	}
	return nil
}

func (a *AuthService) Authorize(ctx context.Context, sessionId string, permission string) (*auth.Session, transport.HError) {
	session := a.sessionRepo.GetSessionUniqueById(sessionId)
	err := verifySession(session)
	if err != nil {
		return nil, err
	}
	for _, p := range session.RolePermissions {
		ok := auth.IsPermissionGrand(permission, p.Permission)
		if ok {
			return session, nil
		}
	}
	return nil, transport.NewForbiddenError(permission)
}

func (a *AuthService) VerifySession(ctx context.Context, sessionId string) (*auth.Session, transport.HError) {
	session := a.sessionRepo.GetSessionUniqueById(sessionId)
	err := verifySession(session)
	if err != nil {
		return nil, err
	}
	return session, nil
}

func verifySession(session *auth.Session) transport.HError {
	if session == nil || session.EmailPasswd == nil {
		return transport.NewPreconditionRequireError(
			transport.PreconditionRequireMsgLogin,
			transport.RequireActionLogin,
		)
	}
	if session.EmailPasswd.Status == auth.EPStatusWaitingVerify {
		return transport.NewPreconditionRequireError(
			transport.PreconditionRequireMsgVerifyEmail,
			transport.RequireActionVerifyEmail,
		)
	}
	if session.User == nil || session.User.Status == auth.UserStatusInitial {
		return transport.NewPreconditionRequireError(
			transport.PreconditionRequireMsgCreateUser,
			transport.RequireActionCreateUser,
		)
	}
	return nil
}
