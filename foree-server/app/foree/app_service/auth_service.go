package foree_service

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"xue.io/go-pay/app/foree/account"
	foree_auth "xue.io/go-pay/app/foree/auth"
	foree_constant "xue.io/go-pay/app/foree/constant"
	foree_logger "xue.io/go-pay/app/foree/logger"
	"xue.io/go-pay/app/foree/referral"
	"xue.io/go-pay/auth"
	"xue.io/go-pay/constant"
	"xue.io/go-pay/server/transport"
	http_util "xue.io/go-pay/util/http"
)

const maxLoginAttempts = 4
const verifyCodeExpiry = 4 * time.Minute
const retrieveTokenExpiry = 5 * time.Minute
const forgetPasswdUpdateInterval = 1 * time.Hour

func NewAuthService(
	db *sql.DB,
	sessionRepo *auth.SessionRepo,
	userRepo *auth.UserRepo,
	emailPasswordRepo *auth.EmailPasswdRepo,
	rolePermissionRepo *auth.RolePermissionRepo,
	userIdentificationRepo *foree_auth.UserIdentificationRepo,
	interacAccountRepo *account.InteracAccountRepo,
	userGroupRepo *auth.UserGroupRepo,
	userExtraRepo *foree_auth.UserExtraRepo,
	userSettingRepo *auth.UserSettingRepo,
	referralRepo *referral.ReferralRepo,
	promotionService *PromotionService,
) *AuthService {
	return &AuthService{
		db:                     db,
		sessionRepo:            sessionRepo,
		userRepo:               userRepo,
		emailPasswordRepo:      emailPasswordRepo,
		rolePermissionRepo:     rolePermissionRepo,
		userIdentificationRepo: userIdentificationRepo,
		interacAccountRepo:     interacAccountRepo,
		userGroupRepo:          userGroupRepo,
		userExtraRepo:          userExtraRepo,
		userSettingRepo:        userSettingRepo,
		referralRepo:           referralRepo,
		promotionService:       promotionService,
	}
}

type AuthService struct {
	db                     *sql.DB
	sessionRepo            *auth.SessionRepo
	userRepo               *auth.UserRepo
	emailPasswordRepo      *auth.EmailPasswdRepo
	rolePermissionRepo     *auth.RolePermissionRepo
	userIdentificationRepo *foree_auth.UserIdentificationRepo
	interacAccountRepo     *account.InteracAccountRepo
	userGroupRepo          *auth.UserGroupRepo
	userExtraRepo          *foree_auth.UserExtraRepo
	userSettingRepo        *auth.UserSettingRepo
	referralRepo           *referral.ReferralRepo
	promotionService       *PromotionService
}

func (a *AuthService) SignUp(ctx context.Context, req SignUpReq) (*UserDTO, transport.HError) {
	// Check if email already exists.
	oldEmail, err := a.emailPasswordRepo.GetUniqueEmailPasswdByEmail(req.Email)
	if err != nil {
		foree_logger.Logger.Error("Signup_FAIL", "ip", loadRealIp(ctx), "email", req.Email, "cause", err.Error())
		return nil, transport.WrapInteralServerError(err)
	}

	if oldEmail != nil {
		foree_logger.Logger.Warn("Signup_FAIL", "ip", loadRealIp(ctx), "email", req.Email, "cause", "email already exists")
		return nil, transport.NewFormError("Invaild signup", "email", "account already exists")
	}

	// Hashing password.
	hashedPasswd, err := auth.HashPassword(req.Password)
	if err != nil {
		foree_logger.Logger.Error("Signup_FAIL", "ip", loadRealIp(ctx), "email", req.Email, "cause", err.Error())
		return nil, transport.WrapInteralServerError(err)
	}

	// Start DB transaction
	dTx, err := a.db.Begin()
	if err != nil {
		foree_logger.Logger.Error("Signup_FAIL", "ip", loadRealIp(ctx), "email", req.Email, "cause", err.Error())
		dTx.Rollback()
		return nil, transport.WrapInteralServerError(err)
	}
	ctx = context.WithValue(ctx, constant.CKdatabaseTransaction, dTx)

	// Create User
	userId, err := a.userRepo.InsertUser(ctx, auth.User{
		Status: auth.UserStatusInitial,
		Email:  req.Email,
	})

	if err != nil {
		dTx.Rollback()
		foree_logger.Logger.Error("Signup_FAIL", "ip", loadRealIp(ctx), "email", req.Email, "cause", err.Error())
		return nil, transport.WrapInteralServerError(err)
	}

	// Create EmailPasswd
	verifyCodeExpiredAt := time.Now().Add(verifyCodeExpiry)
	epId, err := a.emailPasswordRepo.InsertEmailPasswd(ctx, auth.EmailPasswd{
		Email:               req.Email,
		Username:            req.Email,
		Passwd:              hashedPasswd,
		Status:              auth.EPStatusWaitingVerify,
		VerifyCode:          auth.GenerateVerifyCode(),
		VerifyCodeExpiredAt: &verifyCodeExpiredAt,
		OwnerId:             userId,
	})

	if err != nil {
		dTx.Rollback()
		foree_logger.Logger.Error("Signup_FAIL", "ip", loadRealIp(ctx), "email", req.Email, "cause", err.Error())
		return nil, transport.WrapInteralServerError(err)
	}

	//Create userGroup
	_, err = a.userGroupRepo.InsertUserGroup(ctx, auth.UserGroup{
		RoleGroup:             foree_constant.DefaultRoleGroup,
		TransactionLimitGroup: foree_constant.DefaultTransactionLimitGroup,
		FeeGroup:              foree_constant.DefaultFeeGroup,
		OwnerId:               userId,
	})
	if err != nil {
		dTx.Rollback()
		foree_logger.Logger.Error("Signup_FAIL", "ip", loadRealIp(ctx), "userId", userId, "cause", err.Error())
		return nil, transport.WrapInteralServerError(err)
	}

	if err = dTx.Commit(); err != nil {
		foree_logger.Logger.Error("Signup_FAIL", "ip", loadRealIp(ctx), "email", req.Email, "cause", err.Error())
		return nil, transport.WrapInteralServerError(err)
	}

	foree_logger.Logger.Info("Signup_SUCCESS", "ip", loadRealIp(ctx), "userId", userId, "emailPasswordId", epId)

	user, err := a.userRepo.GetUniqueUserById(userId)
	if err != nil {
		foree_logger.Logger.Error("Signup_FAIL", "ip", loadRealIp(ctx), "email", req.Email, "cause", err.Error())
		return nil, transport.WrapInteralServerError(err)
	}

	if user == nil {
		foree_logger.Logger.Error("Signup_FAIL", "ip", loadRealIp(ctx), "email", req.Email, "userId", userId, "cause", "unable to get user")
		return nil, transport.NewInteralServerError("unable to get user with id: `%v`", userId)
	}

	ep, err := a.emailPasswordRepo.GetUniqueEmailPasswdById(epId)

	if err != nil {
		foree_logger.Logger.Error("Signup_FAIL", "ip", loadRealIp(ctx), "email", req.Email, "cause", err.Error())
		return nil, transport.WrapInteralServerError(err)
	}

	if ep == nil {
		foree_logger.Logger.Error("Signup_FAIL", "ip", loadRealIp(ctx), "email", req.Email, "emailPasswordId", epId, "cause", "unable to get created emailPassword")
		return nil, transport.NewInteralServerError("unable to get EmailPasswd with id: `%v`", epId)
	}

	go a.linkReferer(*user, req)

	//TODO: send email. by goroutine

	sessionId, err := a.sessionRepo.InsertSession(auth.Session{
		UserId:      user.ID,
		User:        user,
		EmailPasswd: ep,
	})
	if err != nil {
		foree_logger.Logger.Error("Signup_FAIL", "ip", loadRealIp(ctx), "email", req.Email, "cause", err.Error())
		return nil, transport.WrapInteralServerError(err)
	}

	session := a.sessionRepo.GetSessionUniqueById(sessionId)
	if session == nil {
		foree_logger.Logger.Error("Signup_FAIL", "ip", loadRealIp(ctx), "email", req.Email, "cause", "unable to get created session")
		return nil, transport.NewInteralServerError("sesson `%s` not found", sessionId)
	}
	return NewUserDTO(session), nil
}

func (a *AuthService) allowVerifyEmail(sessionId string) (*auth.Session, transport.HError) {
	session := a.sessionRepo.GetSessionUniqueById(sessionId)
	if session == nil {
		return session, transport.NewUnauthorizedRequestError()
	}
	if session.EmailPasswd.Status != auth.EPStatusWaitingVerify {
		return session, transport.NewPreconditionRequireError(
			transport.PreconditionRequireMsgToMain,
			transport.RequireActionToMain,
		)
	}

	return session, nil
}

func (a *AuthService) VerifyEmail(ctx context.Context, req VerifyEmailReq) (*UserDTO, transport.HError) {
	// Check Allow to VerifyEmail
	session, sErr := a.allowVerifyEmail(req.SessionId)
	if sErr != nil {
		var userId int64
		if session != nil {
			userId = session.UserId
		}
		foree_logger.Logger.Warn("VerifyEmail_FAIL", "ip", loadRealIp(ctx), "userId", userId, "sessionId", req.SessionId, "cause", sErr.Error())
		return nil, sErr
	}

	curEP, err := a.emailPasswordRepo.GetUniqueEmailPasswdById(session.EmailPasswd.ID)

	if err != nil {
		foree_logger.Logger.Warn("VerifyEmail_FAIL", "ip", loadRealIp(ctx), "userId", session.UserId, "cause", err.Error())
		return nil, transport.WrapInteralServerError(err)
	}

	if curEP.VerifyCode != req.Code {
		foree_logger.Logger.Warn("VerifyEmail_FAIL", "ip", loadRealIp(ctx), "userId", session.UserId, "cause", "invalid code")
		return nil, transport.NewFormError("Invalid VerifyEmail Requst", "code", "invalid code")
	}

	if curEP.VerifyCodeExpiredAt == nil || curEP.VerifyCodeExpiredAt.Before(time.Now()) {
		foree_logger.Logger.Warn("VerifyEmail_FAIL", "ip", loadRealIp(ctx), "userId", session.UserId, "cause", "code expired")
		return nil, transport.NewFormError("Invalid VerifyEmail Requst", "code", "code expired")
	}

	// VerifyEmail and update EmailPasswd.
	newEP := *session.EmailPasswd
	newEP.Status = auth.EPStatusActive

	e := a.emailPasswordRepo.UpdateEmailPasswdByEmail(ctx, newEP)
	if e != nil {
		foree_logger.Logger.Warn("VerifyEmail_FAIL", "ip", loadRealIp(ctx), "userId", session.UserId, "cause", e.Error())
		return nil, transport.WrapInteralServerError(e)
	}

	// Update session
	newSession := *session
	newSession.EmailPasswd = &newEP

	session, e = a.sessionRepo.UpdateSession(newSession)
	if e != nil {
		foree_logger.Logger.Warn("VerifyEmail_FAIL", "ip", loadRealIp(ctx), "userId", session.UserId, "cause", e.Error())
		return nil, transport.WrapInteralServerError(e)
	}
	foree_logger.Logger.Info("VerifyEmail_SUCCESS", "ip", loadRealIp(ctx), "userId", session.UserId, "sessionId", req.SessionId)
	return NewUserDTO(session), nil
}

func (a *AuthService) ResendVerifyCode(ctx context.Context, req transport.SessionReq) (*UserDTO, transport.HError) {
	// Check Allow to VerifyEmail
	session, err := a.allowVerifyEmail(req.SessionId)
	if err != nil {
		var userId int64
		if session != nil {
			userId = session.UserId
		}
		foree_logger.Logger.Warn("ResendVerifyCode_FAIL", "ip", loadRealIp(ctx), "userId", userId, "sessionId", req.SessionId, "cause", err.Error())
		return nil, err
	}

	// Change VerifyCode
	verifyCodeExpiredAt := time.Now().Add(verifyCodeExpiry)
	newEP := *session.EmailPasswd
	newEP.VerifyCode = auth.GenerateVerifyCode()
	newEP.VerifyCodeExpiredAt = &verifyCodeExpiredAt

	e := a.emailPasswordRepo.UpdateEmailPasswdByEmail(ctx, newEP)
	if e != nil {
		foree_logger.Logger.Warn("ResendVerifyCode_FAIL", "ip", loadRealIp(ctx), "userId", session.UserId, "cause", e.Error())
		return nil, transport.WrapInteralServerError(e)
	}

	//TODO: send email. by goroutine

	// Update session
	newSession := *session
	newSession.EmailPasswd = &newEP

	_, e = a.sessionRepo.UpdateSession(newSession)
	if e != nil {
		foree_logger.Logger.Warn("ResendVerifyCode_FAIL", "ip", loadRealIp(ctx), "userId", session.UserId, "cause", e.Error())
		return nil, transport.WrapInteralServerError(e)
	}

	foree_logger.Logger.Info("ResendVerifyCode_SUCCESS", "ip", loadRealIp(ctx), "userId", session.UserId, "sessionId", req.SessionId)
	return NewUserDTO(session), nil
}

func (a *AuthService) allowCreateUser(sessionId string) (*auth.Session, transport.HError) {
	session := a.sessionRepo.GetSessionUniqueById(sessionId)
	if session == nil {
		return session, transport.NewUnauthorizedRequestError()
	}

	if session.EmailPasswd.Status == auth.EPStatusWaitingVerify {
		return session, transport.NewPreconditionRequireError(
			transport.PreconditionRequireMsgVerifyEmail,
			transport.RequireActionVerifyEmail,
		)
	}

	if session.User.Status != auth.UserStatusInitial {
		return session, transport.NewPreconditionRequireError(
			transport.PreconditionRequireMsgToMain,
			transport.RequireActionToMain,
		)
	}

	return session, nil
}

func (a *AuthService) CreateUser(ctx context.Context, req CreateUserReq) (*UserDTO, transport.HError) {
	// Check allow to create user
	session, sErr := a.allowCreateUser(req.SessionId)
	if sErr != nil {
		var userId int64
		if session != nil {
			userId = session.UserId
		}
		foree_logger.Logger.Warn("CreateUser_FAIL", "ip", loadRealIp(ctx), "userId", userId, "sessionId", req.SessionId, "cause", sErr.Error())
		return nil, sErr
	}

	curUser, err := a.userRepo.GetUniqueUserById(session.UserId)
	if err != nil {
		foree_logger.Logger.Error("CreateUser_FAIL", "ip", loadRealIp(ctx), "userId", session.UserId, "cause", err.Error())
		return nil, transport.WrapInteralServerError(err)
	}

	// Ask user login
	if curUser.Status != auth.UserStatusInitial {
		a.sessionRepo.Delete(req.SessionId)
		foree_logger.Logger.Warn("CreateUser_FAIL", "ip", loadRealIp(ctx), "userId", session.UserId, "cause", fmt.Sprintf("user in status `%v`", curUser.Status))
		return nil, transport.NewUnauthorizedRequestError()
	}

	dTx, dErr := a.db.Begin()
	if dErr != nil {
		dTx.Rollback()
		foree_logger.Logger.Error("CreateUser_FAIL", "ip", loadRealIp(ctx), "userId", session.UserId, "cause", dErr.Error())
		return nil, transport.WrapInteralServerError(dErr)
	}
	ctx = context.WithValue(ctx, constant.CKdatabaseTransaction, dTx)

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
		foree_logger.Logger.Error("CreateUser_FAIL", "ip", loadRealIp(ctx), "userId", session.UserId, "cause", ier.Error())
		return nil, transport.WrapInteralServerError(ier)
	}

	// Create a new user by updating essential fields.
	newUser := *session.User
	newUser.Status = auth.UserStatusActive
	newUser.FirstName = req.FirstName
	newUser.MiddleName = req.MiddleName
	newUser.LastName = req.LastName
	newUser.Age = req.Age
	newUser.Dob = &req.Dob.Time
	newUser.Address1 = req.Address1
	newUser.Address2 = req.Address2
	newUser.City = req.City
	newUser.Province = req.Province
	newUser.Country = req.Country
	newUser.PostalCode = req.PostalCode
	newUser.PhoneNumber = req.PhoneNumber
	newUser.Email = session.EmailPasswd.Email

	updateErr := a.userRepo.UpdateUserById(ctx, newUser)

	if updateErr != nil {
		dTx.Rollback()
		foree_logger.Logger.Error("CreateUser_FAIL", "userId", session.UserId, "cause", updateErr.Error())
		return nil, transport.WrapInteralServerError(updateErr)
	}

	userExtra := foree_auth.UserExtra{
		Nationality: req.Nationality,
		Pob:         req.Pob,
		OwnerId:     session.UserId,
	}

	_, er := a.userExtraRepo.InsertUserExtra(ctx, userExtra)

	// Create default Interac Account for the user.
	now := time.Now()
	acc := account.InteracAccount{
		FirstName:        req.FirstName,
		MiddleName:       req.MiddleName,
		LastName:         req.LastName,
		Email:            session.EmailPasswd.Email,
		OwnerId:          session.UserId,
		Status:           account.AccountStatusActive,
		LatestActivityAt: &now,
	}
	_, aErr := a.interacAccountRepo.InsertInteracAccount(ctx, acc)
	if aErr != nil {
		dTx.Rollback()
		foree_logger.Logger.Error("Default_Interac_Account_FAIL", "ip", loadRealIp(ctx), "userId", session.UserId, "cause", aErr.Error())
		return nil, transport.WrapInteralServerError(aErr)
	}

	if er != nil {
		dTx.Rollback()
		foree_logger.Logger.Error("CreateUser_FAIL", "ip", loadRealIp(ctx), "userId", session.UserId, "cause", er.Error())
		return nil, transport.WrapInteralServerError(er)
	}

	if err := dTx.Commit(); err != nil {
		foree_logger.Logger.Error("CreateUser_FAIL", "ip", loadRealIp(ctx), "userId", session.UserId, "cause", err.Error())
		return nil, transport.WrapInteralServerError(err)
	}

	foree_logger.Logger.Info("CreateUser_SUCCESS", "ip", loadRealIp(ctx), "userId", session.UserId, "sessionId", req.SessionId)

	userGroup, er := a.userGroupRepo.GetUniqueUserGroupByOwnerId(newUser.ID)
	if er != nil {
		foree_logger.Logger.Error("CreateUser_FAIL", "ip", loadRealIp(ctx), "userId", session.UserId, "cause", er.Error())
		return nil, transport.WrapInteralServerError(er)
	}

	if userGroup == nil {
		foree_logger.Logger.Error("Login_FAIL", "ip", loadRealIp(ctx), "userId", session.UserId, "cause", fmt.Sprintf("userGroup not found for owner `%v`", newUser.ID))
		return nil, transport.NewInteralServerError("userGroup not found for owner `%v`", newUser.ID)
	}

	// Get Permission.
	rolePermissions, pErr := a.rolePermissionRepo.GetAllEnabledRolePermissionByRoleName(userGroup.RoleGroup)
	if pErr != nil {
		foree_logger.Logger.Error("CreateUser_FAIL", "ip", loadRealIp(ctx), "userId", session.UserId, "cause", pErr.Error())
		return nil, transport.WrapInteralServerError(pErr)
	}

	// Update session.
	newSession := *session
	newSession.User = &newUser
	newSession.UserId = newUser.ID
	newSession.RolePermissions = rolePermissions
	newSession.UserGroup = userGroup

	updateSession, sessionErr := a.sessionRepo.UpdateSession(newSession)
	if sessionErr != nil {
		foree_logger.Logger.Error("CreateUser_FAIL", "ip", loadRealIp(ctx), "userId", session.UserId, "cause", sessionErr.Error())
		return nil, transport.WrapInteralServerError(sessionErr)
	}

	go a.promotionService.rewardOnboard(newUser)
	go a.createUserSetting(newUser.ID)

	return NewUserDTO(updateSession), nil
}

// TODO: Login protection on peak volume.
func (a *AuthService) Login(ctx context.Context, req LoginReq) (*UserDTO, transport.HError) {
	// Delete previous token if exists.
	a.sessionRepo.Delete(req.SessionId)

	// Verify email and password
	ep, err := a.emailPasswordRepo.GetUniqueEmailPasswdByEmail(req.Email)

	if err != nil {
		foree_logger.Logger.Error("Login_FAIL", "ip", loadRealIp(ctx), "email", req.Email, "cause", err.Error())
		return nil, transport.WrapInteralServerError(err)
	}

	if ep == nil {
		foree_logger.Logger.Warn("Login_FAIL", "ip", loadRealIp(ctx), "email", req.Email, "cause", fmt.Sprintf("email `%s` not found", req.Email))
		return nil, transport.NewFormError("Invaild signup", "email", "invalid email")
	}

	if ep.Status == auth.EPStatusDelete {
		foree_logger.Logger.Warn("Login_FAIL", "ip", loadRealIp(ctx), "email", req.Email, "cause", fmt.Sprintf("email `%s` is in `%s` status", ep.Email, ep.Status))
		return nil, transport.NewFormError("Invalid login request", "email", "invalid email")
	}

	if ep.Status == auth.EPStatusSuspend {
		foree_logger.Logger.Warn("Login_FAIL", "ip", loadRealIp(ctx), "email", req.Email, "cause", fmt.Sprintf("email `%s` is in `%s` status", ep.Email, ep.Status))
		return nil, transport.NewFormError("Invalid login request", "email", "your account is suspend. please contact us.")
	}

	if ep.LoginAttempts > maxLoginAttempts {
		foree_logger.Logger.Warn("Login_FAIL", "ip", loadRealIp(ctx), "email", req.Email, "cause", fmt.Sprintf("email `%s` has try `%v` times", ep.Email, ep.LoginAttempts))
		return nil, transport.NewFormError("Invalid login request", "password", "max login attempts reached. please contact us.")
	}

	ok := auth.ComparePasswords(req.Password, []byte(ep.Passwd))
	if !ok {
		foree_logger.Logger.Warn("Login_FAIL", "ip", loadRealIp(ctx), "email", req.Email, "cause", "invalid password")
		go func() {
			newEP := *ep
			newEP.LoginAttempts += 1
			if err := a.emailPasswordRepo.UpdateEmailPasswdByEmail(ctx, newEP); err != nil {
				foree_logger.Logger.Error("Login_Attempts_Update", "email", req.Email, "cause", err.Error())
			}
		}()
		return nil, transport.NewFormError("Invaild signup", "password", "Invalid password")
	}

	// Load user(user must exist, but not necessary to be active)
	user, err := a.userRepo.GetUniqueUserById(ep.OwnerId)
	if err != nil {
		foree_logger.Logger.Error("Login_FAIL", "ip", loadRealIp(ctx), "email", req.Email, "cause", err.Error())
		return nil, transport.WrapInteralServerError(err)
	}
	// User must exists
	if user == nil {
		foree_logger.Logger.Error("Login_FAIL", "ip", loadRealIp(ctx), "email", req.Email, "cause", fmt.Sprintf("owner `%v` no found", ep.OwnerId))
		return nil, transport.NewInteralServerError("User `%v` do not exists", ep.OwnerId)
	}

	userGroup, er := a.userGroupRepo.GetUniqueUserGroupByOwnerId(user.ID)
	if er != nil {
		foree_logger.Logger.Error("Login_FAIL", "ip", loadRealIp(ctx), "email", req.Email, "user", user.ID, "cause", er.Error())
		return nil, transport.WrapInteralServerError(er)
	}
	//User group must exists
	if userGroup == nil {
		foree_logger.Logger.Error("Login_FAIL", "ip", loadRealIp(ctx), "email", req.Email, "cause", fmt.Sprintf("userGroup not found with owner `%v`", ep.OwnerId))
		return nil, transport.NewInteralServerError("userGroup not found with owner `%v`", ep.OwnerId)
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

	newSession.Ip = loadRealIp(ctx)
	newSession.UserAgent = loadUserAgent(ctx)

	sessionId, err := a.sessionRepo.InsertSession(newSession)
	if err != nil {
		foree_logger.Logger.Error("Login_FAIL", "ip", loadRealIp(ctx), "email", req.Email, "cause", "fail to insert session to session repo")
		return nil, transport.WrapInteralServerError(err)
	}
	session := a.sessionRepo.GetSessionUniqueById(sessionId)
	if session == nil {
		foree_logger.Logger.Error("Login_FAIL", "ip", loadRealIp(ctx), "email", req.Email, "cause", "fail to get session from session repo")
		return nil, transport.NewInteralServerError("sesson `%s` not found", sessionId)
	}

	go func() {
		newEP := *ep
		newEP.LoginAttempts = 0
		if err := a.emailPasswordRepo.UpdateEmailPasswdByEmail(ctx, newEP); err != nil {
			foree_logger.Logger.Error("Login_reset_login_attempts_FAIL", "emailPasswdId", newEP.ID, "cause", err.Error())
		}

	}()
	foree_logger.Logger.Info("Login_SUCCESS", "ip", loadRealIp(ctx), "email", req.Email, "userAgent", loadUserAgent(ctx), "userId", user.ID)
	return NewUserDTO(session), nil
}

func (a *AuthService) ForgetPasswd(ctx context.Context, req ForgetPasswdReq) (any, transport.HError) {
	// Verify email and password
	ep, err := a.emailPasswordRepo.GetUniqueEmailPasswdByEmail(req.Email)

	if err != nil {
		foree_logger.Logger.Error("ForgetPasswd_FAIL", "ip", loadRealIp(ctx), "email", req.Email, "cause", err.Error())
		return nil, transport.WrapInteralServerError(err)

	}

	if ep == nil {
		foree_logger.Logger.Warn("ForgetPasswd_FAIL", "ip", loadRealIp(ctx), "email", req.Email, "cause", fmt.Sprintf("email `%s` not found", req.Email))
		return nil, transport.NewFormError("Invaild forget password request", "email", "invalid email")
	}

	if ep.LatestForgetPasswdUpdatedAt != nil && ep.LatestForgetPasswdUpdatedAt.After(time.Now().Add(-1*forgetPasswdUpdateInterval)) {
		foree_logger.Logger.Warn("ForgetPasswd_FAIL", "ip", loadRealIp(ctx), "email", req.Email, "cause", "forget password update is to frequent")
		return nil, transport.NewFormError("Invaild forget password request", "email", "password update too frequent")
	}

	newEP := *ep
	retrieveTokenExpiredAt := time.Now().Add(retrieveTokenExpiry)
	newEP.RetrieveToken = auth.GenerateVerifyCode()
	newEP.RetrieveTokenExpiredAt = &retrieveTokenExpiredAt

	err = a.emailPasswordRepo.UpdateEmailPasswdByEmail(ctx, newEP)
	if err != nil {
		foree_logger.Logger.Error("ForgetPasswd_FAIL", "ip", loadRealIp(ctx), "email", req.Email, "cause", err.Error())
		transport.WrapInteralServerError(err)
	}

	foree_logger.Logger.Info("ForgetPasswd_SUCCESS", "ip", loadRealIp(ctx), "email", req.Email)

	//TODO: send email with code

	return nil, nil
}

func (a *AuthService) ForgetPasswdVerify(ctx context.Context, req ForgetPasswdVerifyReq) (any, transport.HError) {
	ep, err := a.emailPasswordRepo.GetUniqueEmailPasswdByEmail(req.Email)

	if err != nil {
		foree_logger.Logger.Error("ForgetPasswdVerify_FAIL", "ip", loadRealIp(ctx), "email", req.Email, "cause", err.Error())
		return nil, transport.WrapInteralServerError(err)

	}

	if ep.LatestForgetPasswdUpdatedAt != nil && ep.LatestForgetPasswdUpdatedAt.After(time.Now().Add(-1*forgetPasswdUpdateInterval)) {
		foree_logger.Logger.Warn("ForgetPasswdVerify_FAIL", "ip", loadRealIp(ctx), "email", req.Email, "cause", "forget password update is to frequent")
		return nil, transport.NewFormError("Invaild forget password request", "email", "password update too frequent")
	}

	if ep == nil {
		foree_logger.Logger.Warn("ForgetPasswdVerify_FAIL", "ip", loadRealIp(ctx), "email", req.Email, "cause", fmt.Sprintf("email `%s` not found", req.Email))
		return nil, transport.NewFormError("Invaild forget password", "email", "invalid email")
	}

	if ep.RetrieveTokenExpiredAt == nil || ep.RetrieveTokenExpiredAt.Before(time.Now()) {
		foree_logger.Logger.Warn("ForgetPasswdVerify_FAIL", "ip", loadRealIp(ctx), "email", req.Email, "cause", "retrieve code expire")
		return nil, transport.NewFormError("Invaild forget password", "retrieveCode", "retrieve code expired")
	}

	if ep.RetrieveToken != req.RetrieveCode {
		foree_logger.Logger.Warn("ForgetPasswdVerify_FAIL", "ip", loadRealIp(ctx), "email", req.Email, "cause", "invalid retrieve code")
		return nil, transport.NewFormError("Invaild forget password", "retrieveCode", "invalid retrieve code")
	}

	foree_logger.Logger.Info("ForgetPasswdVerify_SUCCESS", "ip", loadRealIp(ctx), "email", req.Email)

	return nil, nil
}

func (a *AuthService) ForgetPasswdUpdate(ctx context.Context, req ForgetPasswdUpdateReq) (any, transport.HError) {
	ep, err := a.emailPasswordRepo.GetUniqueEmailPasswdByEmail(req.Email)

	if err != nil {
		foree_logger.Logger.Error("ForgetPasswdUpdate_FAIL", "ip", loadRealIp(ctx), "email", req.Email, "cause", err.Error())
		return nil, transport.WrapInteralServerError(err)

	}

	if ep.LatestForgetPasswdUpdatedAt != nil && ep.LatestForgetPasswdUpdatedAt.After(time.Now().Add(-1*forgetPasswdUpdateInterval)) {
		foree_logger.Logger.Warn("ForgetPasswdUpdate_FAIL", "ip", loadRealIp(ctx), "email", req.Email, "cause", "forget password update is to frequent")
		return nil, transport.NewFormError("Invaild forget password request", "email", "password update too frequent")
	}

	if ep == nil {
		foree_logger.Logger.Warn("ForgetPasswdUpdate_FAIL", "ip", loadRealIp(ctx), "email", req.Email, "cause", fmt.Sprintf("email `%s` not found", req.Email))
		return nil, transport.NewFormError("Invaild forget password", "email", "invalid email")
	}

	if ep.RetrieveTokenExpiredAt == nil || ep.RetrieveTokenExpiredAt.Before(time.Now()) {
		foree_logger.Logger.Warn("ForgetPasswdUpdate_FAIL", "ip", loadRealIp(ctx), "email", req.Email, "cause", "retrieve code expire")
		return nil, transport.NewFormError("Invaild forget password", "retrieveCode", "retrieve code expired")
	}

	if ep.RetrieveToken != req.RetrieveCode {
		foree_logger.Logger.Warn("ForgetPasswdUpdate_FAIL", "ip", loadRealIp(ctx), "email", req.Email, "cause", "invalid retrieve code")
		return nil, transport.NewFormError("Invaild forget password", "retrieveCode", "invalid retrieve code")
	}

	hashedPasswd, err := auth.HashPassword(req.NewPassword)
	if err != nil {
		foree_logger.Logger.Error("ForgetPasswdUpdate_FAIL", "ip", loadRealIp(ctx), "email", req.Email, "cause", err.Error())
		return nil, transport.WrapInteralServerError(err)
	}

	now := time.Now()
	newEP := *ep
	newEP.Passwd = hashedPasswd
	newEP.LatestForgetPasswdUpdatedAt = &now
	newEP.RetrieveTokenExpiredAt = nil

	err = a.emailPasswordRepo.UpdateEmailPasswdByEmail(context.Background(), newEP)

	if err != nil {
		foree_logger.Logger.Error("ForgetPasswdUpdate_FAIL", "ip", loadRealIp(ctx), "email", req.Email, "cause", err.Error())
		transport.WrapInteralServerError(err)
	}

	//TODO: send email to user for notice

	foree_logger.Logger.Info("ForgetPasswdUpdate_SUCCESS", "ip", loadRealIp(ctx), "email", req.Email)
	return nil, nil
}

func (a *AuthService) Logout(ctx context.Context, req transport.SessionReq) (*auth.Session, transport.HError) {
	oldSession, _ := a.VerifySession(ctx, req.SessionId)
	a.sessionRepo.Delete(req.SessionId)
	if oldSession != nil {
		foree_logger.Logger.Debug("Logout_SUCCESS", "ip", loadRealIp(ctx), "userId", oldSession.UserId, "sessionId", req.SessionId)
	} else {
		foree_logger.Logger.Debug("Logout_SUCCESS", "ip", loadRealIp(ctx), "sessionId", req.SessionId)
	}
	return nil, nil
}

func (a *AuthService) GetUser(ctx context.Context, req transport.SessionReq) (*UserDTO, transport.HError) {
	session, sErr := a.VerifySession(ctx, req.SessionId)
	if sErr != nil {
		var userId int64
		if session != nil {
			userId = session.UserId
		}
		// Normal error when the token expired
		foree_logger.Logger.Info("GetUser_FAIL", "ip", loadRealIp(ctx), "userId", userId, "sessionId", req.SessionId, "cause", sErr.Error())
		return nil, sErr
	}
	foree_logger.Logger.Info("GetUser_SUCCESS", "ip", loadRealIp(ctx), "userId", session.UserId, "sessionId", req.SessionId)
	return NewUserDTO(session), nil
}

func (a *AuthService) GetUserDetail(ctx context.Context, req transport.SessionReq) (any, transport.HError) {
	session, sErr := a.VerifySession(ctx, req.SessionId)
	if sErr != nil {
		var userId int64
		if session != nil {
			userId = session.UserId
		}
		// Normal error when the token expired
		foree_logger.Logger.Info("GetUserDetail_FAIL", "ip", loadRealIp(ctx), "userId", userId, "sessionId", req.SessionId, "cause", sErr.Error())
		return nil, sErr
	}

	user, err := a.userRepo.GetUniqueUserById(session.UserId)
	if err != nil {
		foree_logger.Logger.Error("GetUserDetail_FAIL", "ip", loadRealIp(ctx), "userId", session.UserId, "cause", err.Error())
		return nil, transport.WrapInteralServerError(err)

	}

	if user == nil {
		foree_logger.Logger.Error("GetUserDetail_FAIL", "ip", loadRealIp(ctx), "userId", session.UserId, "emailPasswordId", session.EmailPasswd.ID, "cause", "user no found.")
		return nil, transport.NewInteralServerError("user no found with id `%v`", session.UserId)
	}

	return NewUserDetailDTO(user), nil
}

func (a *AuthService) createUserSetting(ownerId int64) {
	_, err := a.userSettingRepo.InsertUserSetting(context.TODO(), auth.UserSetting{
		IsInAppNotificationEnable:  true,
		IsPushNotificationEnable:   true,
		IsEmailNotificationsEnable: true,
		OwnerId:                    ownerId,
	})

	if err != nil {
		foree_logger.Logger.Error("createUserSetting_FAIL", "userId", ownerId, "cause", err.Error())
	}

}

func (a *AuthService) GetUserSetting(ctx context.Context, req transport.SessionReq) (*UserSettingDTO, transport.HError) {
	session, sErr := a.VerifySession(ctx, req.SessionId)
	if sErr != nil {
		var userId int64
		if session != nil {
			userId = session.UserId
		}
		// Normal error when the token expired
		foree_logger.Logger.Info("GetUserSetting_FAIL", "ip", loadRealIp(ctx), "userId", userId, "sessionId", req.SessionId, "cause", sErr.Error())
		return nil, sErr
	}

	us, err := a.userSettingRepo.GetUniqueUserSettingByOwnerId(session.UserId)
	if err != nil {
		foree_logger.Logger.Error("GetUserSetting_FAIL", "ip", loadRealIp(ctx), "userId", session.UserId, "cause", err.Error())
		return nil, transport.WrapInteralServerError(err)

	}

	if us != nil {
		return NewUserSettingDTO(us), nil
	}

	// Create one.
	_, err = a.userSettingRepo.InsertUserSetting(context.TODO(), auth.UserSetting{
		IsInAppNotificationEnable:  true,
		IsPushNotificationEnable:   true,
		IsEmailNotificationsEnable: true,
		OwnerId:                    session.UserId,
	})

	if err != nil {
		foree_logger.Logger.Error("GetUserSetting_FAIL", "userId", session.UserId, "cause", err.Error())
		return nil, transport.WrapInteralServerError(err)
	}

	us, err = a.userSettingRepo.GetUniqueUserSettingByOwnerId(session.UserId)
	if err != nil {
		foree_logger.Logger.Error("GetUserSetting_FAIL", "ip", loadRealIp(ctx), "userId", session.UserId, "cause", err.Error())
		return nil, transport.WrapInteralServerError(err)

	}

	if us == nil {
		foree_logger.Logger.Error("GetUserSetting_FAIL", "ip", loadRealIp(ctx), "userId", session.UserId, "cause", "userSetting no found")
		return nil, transport.NewInteralServerError("userSetting no found")
	}
	return NewUserSettingDTO(us), nil
}

func (a *AuthService) UpdateUserSetting(ctx context.Context, req UpdateUserSetting) (any, transport.HError) {
	session, sErr := a.VerifySession(ctx, req.SessionId)
	if sErr != nil {
		var userId int64
		if session != nil {
			userId = session.UserId
		}
		// Normal error when the token expired
		foree_logger.Logger.Info("UpdateUserSetting_FAIL", "ip", loadRealIp(ctx), "userId", userId, "sessionId", req.SessionId, "cause", sErr.Error())
		return nil, sErr
	}

	us, err := a.userSettingRepo.GetUniqueUserSettingByOwnerId(session.UserId)
	if err != nil {
		foree_logger.Logger.Error("UpdateUserSetting_FAIL", "ip", loadRealIp(ctx), "userId", session.UserId, "cause", err.Error())
		return nil, transport.WrapInteralServerError(err)
	}

	if us == nil {
		foree_logger.Logger.Error("UpdateUserSetting_FAIL", "ip", loadRealIp(ctx), "userId", session.UserId, "cause", "useSetting no found")
		return nil, transport.NewInteralServerError("userSetting no found")
	}

	err = a.userSettingRepo.UpdateUserSettingByOwnerId(ctx, auth.UserSetting{
		IsInAppNotificationEnable:  req.IsInAppNotificationEnable,
		IsPushNotificationEnable:   req.IsPushNotificationEnable,
		IsEmailNotificationsEnable: req.IsEmailNotificationsEnable,
		OwnerId:                    session.UserId,
	})

	if err != nil {
		foree_logger.Logger.Error("UpdateUserSetting_FAIL", "ip", loadRealIp(ctx), "userId", session.UserId, "cause", err.Error())
		return nil, transport.WrapInteralServerError(err)
	}
	return nil, nil
}

func (a *AuthService) GetUserExtra(ctx context.Context, req transport.SessionReq) (*UserExtraDTO, transport.HError) {
	session, sErr := a.VerifySession(ctx, req.SessionId)
	if sErr != nil {
		var userId int64
		if session != nil {
			userId = session.UserId
		}
		// Normal error when the token expired
		foree_logger.Logger.Info("GetUserExtra_FAIL", "ip", loadRealIp(ctx), "userId", userId, "sessionId", req.SessionId, "cause", sErr.Error())
		return nil, sErr
	}

	ue, err := a.userExtraRepo.GetUniqueUserExtraByOwnerId(session.UserId)
	if err != nil {
		foree_logger.Logger.Error("GetUserExtra_FAIL", "ip", loadRealIp(ctx), "userId", session.UserId, "cause", err.Error())
		return nil, transport.WrapInteralServerError(err)
	}

	if ue == nil {
		return nil, nil
	}

	return NewUserExtraDTO(ue), nil
}

func (a *AuthService) UpdateUserPhoneNumber(ctx context.Context, req UpdatePhoneNumberReq) (any, transport.HError) {
	session, sErr := a.VerifySession(ctx, req.SessionId)
	if sErr != nil {
		var userId int64
		if session != nil {
			userId = session.UserId
		}
		// Normal error when the token expired
		foree_logger.Logger.Info("UpdateUserPhoneNumber_FAIL", "ip", loadRealIp(ctx), "userId", userId, "sessionId", req.SessionId, "cause", sErr.Error())
		return nil, sErr
	}

	user, err := a.userRepo.GetUniqueUserById(session.UserId)
	if err != nil {
		foree_logger.Logger.Error("UpdateUserPhoneNumber_FAIL", "ip", loadRealIp(ctx), "userId", session.UserId, "cause", err.Error())
		return nil, transport.WrapInteralServerError(err)

	}

	if user == nil {
		foree_logger.Logger.Error("UpdateUserPhoneNumber_FAIL", "ip", loadRealIp(ctx), "userId", session.UserId, "emailPasswordId", session.EmailPasswd.ID, "cause", "user no found.")
		return nil, transport.NewInteralServerError("user no found with id `%v`", session.UserId)
	}

	newUser := *user
	newUser.PhoneNumber = req.PhoneNumber

	err = a.userRepo.UpdateUserById(ctx, newUser)
	if err != nil {
		foree_logger.Logger.Error("UpdateUserPhoneNumber_FAIL", "ip", loadRealIp(ctx), "userId", session.UserId, "cause", err.Error())
		return nil, transport.WrapInteralServerError(err)
	}
	foree_logger.Logger.Info("UpdateUserPhoneNumber_SUCCESS", "ip", loadRealIp(ctx), "userId", session.UserId)
	return nil, nil
}

func (a *AuthService) UpdateUserAddress(ctx context.Context, req UpdateAddressReq) (any, transport.HError) {
	session, sErr := a.VerifySession(ctx, req.SessionId)
	if sErr != nil {
		var userId int64
		if session != nil {
			userId = session.UserId
		}
		// Normal error when the token expired
		foree_logger.Logger.Info("UpdateUserAddress_FAIL", "ip", loadRealIp(ctx), "userId", userId, "sessionId", req.SessionId, "cause", sErr.Error())
		return nil, sErr
	}

	user, err := a.userRepo.GetUniqueUserById(session.UserId)
	if err != nil {
		foree_logger.Logger.Error("UpdateUserAddress_FAIL", "ip", loadRealIp(ctx), "userId", session.UserId, "cause", err.Error())
		return nil, transport.WrapInteralServerError(err)

	}

	if user == nil {
		foree_logger.Logger.Error("UpdateUserAddress_FAIL", "ip", loadRealIp(ctx), "userId", session.UserId, "emailPasswordId", session.EmailPasswd.ID, "cause", "user no found.")
		return nil, transport.NewInteralServerError("user no found with id `%v`", session.UserId)
	}

	newUser := *user
	newUser.Address1 = req.Address1
	newUser.Address2 = req.Address2
	newUser.City = req.City
	newUser.Province = req.Province
	newUser.Country = req.Country
	newUser.PostalCode = req.PostalCode

	err = a.userRepo.UpdateUserById(ctx, newUser)
	if err != nil {
		foree_logger.Logger.Error("UpdateUserAddress_FAIL", "ip", loadRealIp(ctx), "userId", session.UserId, "cause", err.Error())
		return nil, transport.WrapInteralServerError(err)
	}
	foree_logger.Logger.Info("UpdateUserAddress_SUCCESS", "ip", loadRealIp(ctx), "userId", session.UserId, "emailPasswordId", session.EmailPasswd.ID)
	return nil, nil
}

func (a *AuthService) UpdatePasswd(ctx context.Context, req UpdatePasswdReq) (*auth.Session, transport.HError) {
	session, sErr := a.VerifySession(ctx, req.SessionId)
	if sErr != nil {
		var userId int64
		if session != nil {
			userId = session.UserId
		}
		// Normal error when the token expired
		foree_logger.Logger.Info("UpdatePasswd_FAIL", "ip", loadRealIp(ctx), "userId", userId, "sessionId", req.SessionId, "cause", sErr.Error())
		return nil, sErr
	}

	ep, err := a.emailPasswordRepo.GetUniqueEmailPasswdById(session.EmailPasswd.ID)

	if err != nil {
		foree_logger.Logger.Error("UpdatePasswd_FAIL", "ip", loadRealIp(ctx), "userId", session.UserId, "cause", err.Error())
		return nil, transport.WrapInteralServerError(err)
	}

	if ep == nil {
		foree_logger.Logger.Error("UpdatePasswd_FAIL", "ip", loadRealIp(ctx), "userId", session.UserId, "emailPasswordId", session.EmailPasswd.ID, "cause", "emailPassword no found")
		return nil, transport.NewInteralServerError("emailPassword no found with id `%v`", session.EmailPasswd.ID)
	}

	ok := auth.ComparePasswords(req.OldPasswd, []byte(ep.Passwd))
	if !ok {
		foree_logger.Logger.Warn("UpdatePasswd_FAIL", "ip", loadRealIp(ctx), "userId", session.UserId, "cause", "invalid old, password")
		return nil, transport.NewFormError("Invalid change passwd request", "oldPasswd", "Invalid old password")
	}

	hashed, err := auth.HashPassword(req.NewPasswd)
	if err != nil {
		foree_logger.Logger.Error("UpdatePasswd_FAIL", "ip", loadRealIp(ctx), "userId", session.UserId, "sessionId", req.SessionId, "cause", err.Error())
		return nil, transport.WrapInteralServerError(err)
	}
	newEp := *session.EmailPasswd
	newEp.Passwd = hashed

	err = a.emailPasswordRepo.UpdateEmailPasswdByEmail(ctx, newEp)
	if err != nil {
		foree_logger.Logger.Error("UpdatePasswd_FAIL", "ip", loadRealIp(ctx), "userId", session.UserId, "sessionId", req.SessionId, "cause", err.Error())
		return nil, transport.WrapInteralServerError(err)
	}
	foree_logger.Logger.Info("UpdatePasswd_SUCCESS", "ip", loadRealIp(ctx), "userId", session.UserId, "sessionId")
	return nil, nil
}

func (a *AuthService) GetSession(ctx context.Context, sessionId string) (*auth.Session, transport.HError) {
	session := a.sessionRepo.GetSessionUniqueById(sessionId)
	if session == nil {
		return nil, transport.NewUnauthorizedRequestError()
	}
	return session, nil
}

func (a *AuthService) Authorize(ctx context.Context, sessionId string, permission string) (*auth.Session, transport.HError) {
	session := a.sessionRepo.GetSessionUniqueById(sessionId)
	err := verifySession(session)
	if err != nil {
		return session, err
	}
	if session.RolePermissions == nil {
		return session, transport.NewForbiddenError(permission)
	}

	for _, p := range session.RolePermissions {
		if p == nil {
			continue
		}
		ok := auth.IsPermissionGrand(permission, p.Permission)
		if ok {
			return session, nil
		}
	}
	return session, transport.NewForbiddenError(permission)
}

func (a *AuthService) VerifySession(ctx context.Context, sessionId string) (*auth.Session, transport.HError) {
	session := a.sessionRepo.GetSessionUniqueById(sessionId)
	err := verifySession(session)
	if err != nil {
		return session, err
	}
	return session, nil
}

func (a *AuthService) linkReferer(registerUser auth.User, req SignUpReq) {
	if req.ReferrerReference == "" {
		return
	}

	referrer, err := a.userExtraRepo.GetUniqueUserExtraByUserReference(req.ReferrerReference)

	if err != nil {
		foree_logger.Logger.Error("Link_Referer_FAIL", "userId", registerUser.ID, "referrerReference", req.ReferrerReference, "cause", err.Error())
		return
	}
	if referrer == nil {
		foree_logger.Logger.Warn("Link_Referer_FAIL", "userId", registerUser.ID, "referrerReference", req.ReferrerReference, "cause", "unknown referrerReference")
		return
	}

	now := time.Now()
	referral := referral.Referral{
		ReferrerId: referrer.ID,
		RefereeId:  registerUser.ID,
		AcceptAt:   &now,
	}

	referralId, err := a.referralRepo.InsertReferral(referral)
	if err != nil {
		foree_logger.Logger.Error("Link_Referer_FAIL", "userId", registerUser.ID, "referrerReference", req.ReferrerReference, "cause", err.Error())
		return
	}
	foree_logger.Logger.Info("Link_Referer_SUCCESS", "userId", registerUser.ID, "referrerReference", req.ReferrerReference, "referrerId", referralId)
}

func verifySession(session *auth.Session) transport.HError {
	if session == nil {
		return transport.NewUnauthorizedRequestError()
	}
	if session.EmailPasswd.Status == auth.EPStatusWaitingVerify {
		return transport.NewPreconditionRequireError(
			transport.PreconditionRequireMsgVerifyEmail,
			transport.RequireActionVerifyEmail,
		)
	}
	if session.User.Status == auth.UserStatusInitial {
		return transport.NewPreconditionRequireError(
			transport.PreconditionRequireMsgCreateUser,
			transport.RequireActionCreateUser,
		)
	}
	return nil
}

func loadRealIp(ctx context.Context) string {
	req, ok := ctx.Value(constant.CKHttpRequest).(*http.Request)
	if !ok {
		return ""
	}
	return http_util.LoadRealIp(req)
}

func loadUserAgent(ctx context.Context) string {
	req, ok := ctx.Value(constant.CKHttpRequest).(*http.Request)
	if !ok {
		return ""
	}
	return req.Header.Get("User-Agent")
}
