package foree_service

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"sync"
	"time"

	"xue.io/go-pay/app/foree/account"
	foree_auth "xue.io/go-pay/app/foree/auth"
	foree_constant "xue.io/go-pay/app/foree/constant"
	foree_logger "xue.io/go-pay/app/foree/logger"
	"xue.io/go-pay/app/foree/promotion"
	"xue.io/go-pay/app/foree/referral"
	"xue.io/go-pay/app/foree/transaction"
	"xue.io/go-pay/auth"
	"xue.io/go-pay/constant"
	"xue.io/go-pay/server/transport"
	http_util "xue.io/go-pay/util/http"
)

const maxLoginAttempts = 4
const PromotionOnboard = "ONBOARD_PROMOTION"
const PromotionReferral = "REFERRAL_PROMOTION"
const promotionCacheTimeout = 15 * time.Minute
const verifyCodeExpiry = 4 * time.Minute

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
	referralRepo *referral.ReferralRepo,
	rewardRepo *transaction.RewardRepo,
	promotionRepo *promotion.PromotionRepo,
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
		referralRepo:           referralRepo,
		rewardRepo:             rewardRepo,
		promotionRepo:          promotionRepo,
	}
}

type AuthService struct {
	db                       *sql.DB
	sessionRepo              *auth.SessionRepo
	userRepo                 *auth.UserRepo
	emailPasswordRepo        *auth.EmailPasswdRepo
	rolePermissionRepo       *auth.RolePermissionRepo
	userIdentificationRepo   *foree_auth.UserIdentificationRepo
	interacAccountRepo       *account.InteracAccountRepo
	userGroupRepo            *auth.UserGroupRepo
	userExtraRepo            *foree_auth.UserExtraRepo
	referralRepo             *referral.ReferralRepo
	rewardRepo               *transaction.RewardRepo
	promotionRepo            *promotion.PromotionRepo
	promotionCache           map[string]CacheItem[promotion.Promotion]
	promotionCacheRWLock     sync.RWMutex
	promotionCacheUpdateLock sync.RWMutex
}

func (a *AuthService) SignUp(ctx context.Context, req SignUpReq) (*UserDTO, transport.HError) {
	// Check if email already exists.
	oldEmail, err := a.emailPasswordRepo.GetUniqueEmailPasswdByEmail(req.Email)
	if err != nil {
		foree_logger.Logger.Error("Signup_Fail", "ip", loadRealIp(ctx), "email", req.Email, "cause", err.Error())
		return nil, transport.WrapInteralServerError(err)
	}

	if oldEmail != nil {
		foree_logger.Logger.Warn("Signup_Fail", "ip", loadRealIp(ctx), "email", req.Email, "cause", "email already exists")
		return nil, transport.NewFormError("Invaild signup", "email", "invalid email")
	}

	// Hashing password.
	hashedPasswd, err := auth.HashPassword(req.Password)
	if err != nil {
		foree_logger.Logger.Error("Signup_Fail", "ip", loadRealIp(ctx), "email", req.Email, "cause", err.Error())
		return nil, transport.WrapInteralServerError(err)
	}

	// Start DB transaction
	dTx, err := a.db.Begin()
	if err != nil {
		foree_logger.Logger.Error("Signup_Fail", "ip", loadRealIp(ctx), "email", req.Email, "cause", err.Error())
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
		foree_logger.Logger.Error("Signup_Fail", "ip", loadRealIp(ctx), "email", req.Email, "cause", err.Error())
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
		foree_logger.Logger.Error("Signup_Fail", "ip", loadRealIp(ctx), "email", req.Email, "cause", err.Error())
		return nil, transport.WrapInteralServerError(err)
	}

	if err = dTx.Commit(); err != nil {
		foree_logger.Logger.Error("Signup_Fail", "ip", loadRealIp(ctx), "email", req.Email, "cause", err.Error())
		return nil, transport.WrapInteralServerError(err)
	}

	foree_logger.Logger.Info("Signup_Success", "ip", loadRealIp(ctx), "userId", userId, "emailPasswordId", epId)

	user, err := a.userRepo.GetUniqueUserById(userId)
	if err != nil {
		foree_logger.Logger.Error("Signup_Fail", "ip", loadRealIp(ctx), "email", req.Email, "cause", err.Error())
		return nil, transport.WrapInteralServerError(err)
	}

	if user == nil {
		foree_logger.Logger.Error("Signup_Fail", "ip", loadRealIp(ctx), "email", req.Email, "userId", userId, "cause", "unable to get user")
		return nil, transport.NewInteralServerError("unable to get user with id: `%v`", userId)
	}

	ep, err := a.emailPasswordRepo.GetUniqueEmailPasswdById(epId)

	if err != nil {
		foree_logger.Logger.Error("Signup_Fail", "ip", loadRealIp(ctx), "email", req.Email, "cause", err.Error())
		return nil, transport.WrapInteralServerError(err)
	}

	if ep == nil {
		foree_logger.Logger.Error("Signup_Fail", "ip", loadRealIp(ctx), "email", req.Email, "emailPasswordId", epId, "cause", "unable to get created emailPassword")
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
		foree_logger.Logger.Error("Signup_Fail", "ip", loadRealIp(ctx), "email", req.Email, "cause", err.Error())
		return nil, transport.WrapInteralServerError(err)
	}

	session := a.sessionRepo.GetSessionUniqueById(sessionId)
	if session == nil {
		foree_logger.Logger.Error("Signup_Fail", "ip", loadRealIp(ctx), "email", req.Email, "cause", "unable to get created session")
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
		foree_logger.Logger.Warn("VerifyEmail_Fail", "ip", loadRealIp(ctx), "userId", userId, "cause", sErr.Error())
		return nil, sErr
	}

	curEP, err := a.emailPasswordRepo.GetUniqueEmailPasswdById(session.EmailPasswd.ID)

	if err != nil {
		foree_logger.Logger.Warn("VerifyEmail_Fail", "ip", loadRealIp(ctx), "userId", session.UserId, "cause", err.Error())
		return nil, transport.WrapInteralServerError(err)
	}

	if curEP.VerifyCode != req.Code {
		foree_logger.Logger.Warn("VerifyEmail_Fail", "ip", loadRealIp(ctx), "userId", session.UserId, "cause", "invalid code")
		return nil, transport.NewFormError("Invalid VerifyEmail Requst", "code", "invalid code")
	}

	if curEP.VerifyCodeExpiredAt == nil || curEP.VerifyCodeExpiredAt.Before(time.Now()) {
		foree_logger.Logger.Warn("VerifyEmail_Fail", "ip", loadRealIp(ctx), "userId", session.UserId, "cause", "code expired")
		return nil, transport.NewFormError("Invalid VerifyEmail Requst", "code", "code expired")
	}

	// VerifyEmail and update EmailPasswd.
	newEP := *session.EmailPasswd
	newEP.Status = auth.EPStatusActive

	e := a.emailPasswordRepo.UpdateEmailPasswdByEmail(ctx, newEP)
	if e != nil {
		foree_logger.Logger.Warn("VerifyEmail_Fail", "ip", loadRealIp(ctx), "userId", session.UserId, "cause", e.Error())
		return nil, transport.WrapInteralServerError(e)
	}

	// Update session
	newSession := *session
	newSession.EmailPasswd = &newEP

	session, e = a.sessionRepo.UpdateSession(newSession)
	if e != nil {
		foree_logger.Logger.Warn("VerifyEmail_Fail", "ip", loadRealIp(ctx), "userId", session.UserId, "cause", e.Error())
		return nil, transport.WrapInteralServerError(e)
	}
	foree_logger.Logger.Info("VerifyEmail_Success", "ip", loadRealIp(ctx), "userId", session.UserId)
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
		foree_logger.Logger.Warn("ResendVerifyCode_Fail", "ip", loadRealIp(ctx), "userId", userId, "cause", err.Error())
		return nil, err
	}

	// Change VerifyCode
	verifyCodeExpiredAt := time.Now().Add(verifyCodeExpiry)
	newEP := *session.EmailPasswd
	newEP.VerifyCode = auth.GenerateVerifyCode()
	newEP.VerifyCodeExpiredAt = &verifyCodeExpiredAt

	e := a.emailPasswordRepo.UpdateEmailPasswdByEmail(ctx, newEP)
	if e != nil {
		foree_logger.Logger.Warn("ResendVerifyCode_Fail", "ip", loadRealIp(ctx), "userId", session.UserId, "cause", e.Error())
		return nil, transport.WrapInteralServerError(e)
	}

	//TODO: send email. by goroutine

	// Update session
	newSession := *session
	newSession.EmailPasswd = &newEP

	_, e = a.sessionRepo.UpdateSession(newSession)
	if e != nil {
		foree_logger.Logger.Warn("ResendVerifyCode_Fail", "ip", loadRealIp(ctx), "userId", session.UserId, "cause", e.Error())
		return nil, transport.WrapInteralServerError(e)
	}

	foree_logger.Logger.Info("ResendVerifyCode_Success", "ip", loadRealIp(ctx), "userId", session.UserId)
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
		foree_logger.Logger.Warn("CreateUser_Fail", "ip", loadRealIp(ctx), "userId", userId, "cause", sErr.Error())
		return nil, sErr
	}

	curUser, err := a.userRepo.GetUniqueUserById(session.UserId)
	if err != nil {
		foree_logger.Logger.Error("CreateUser_Fail", "ip", loadRealIp(ctx), "userId", session.UserId, "cause", err.Error())
		return nil, transport.WrapInteralServerError(err)
	}

	// Ask user login
	if curUser.Status != auth.UserStatusInitial {
		a.sessionRepo.Delete(req.SessionId)
		foree_logger.Logger.Warn("CreateUser_Fail", "ip", loadRealIp(ctx), "userId", session.UserId, "cause", fmt.Sprintf("user in status `%v`", curUser.Status))
		return nil, transport.NewUnauthorizedRequestError()
	}

	dTx, dErr := a.db.Begin()
	if dErr != nil {
		dTx.Rollback()
		foree_logger.Logger.Error("CreateUser_Fail", "ip", loadRealIp(ctx), "userId", session.UserId, "cause", dErr.Error())
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
		foree_logger.Logger.Error("CreateUser_Fail", "ip", loadRealIp(ctx), "userId", session.UserId, "cause", ier.Error())
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
		foree_logger.Logger.Error("CreateUser_Fail", "userId", session.UserId, "cause", updateErr.Error())
		return nil, transport.WrapInteralServerError(updateErr)
	}

	userExtra := foree_auth.UserExtra{
		Nationality: req.Nationality,
		Pob:         req.Pob,
		OwnerId:     session.UserId,
	}

	_, er := a.userExtraRepo.InsertUserExtra(ctx, userExtra)
	if er != nil {
		dTx.Rollback()
		foree_logger.Logger.Error("CreateUser_Fail", "ip", loadRealIp(ctx), "userId", session.UserId, "cause", er.Error())
		return nil, transport.WrapInteralServerError(er)
	}

	//Create userGroup
	_, er = a.userGroupRepo.InsertUserGroup(ctx, auth.UserGroup{
		RoleGroup:             foree_constant.DefaultRoleGroup,
		TransactionLimitGroup: foree_constant.DefaultTransactionLimitGroup,
		OwnerId:               newUser.ID,
	})
	if er != nil {
		dTx.Rollback()
		foree_logger.Logger.Error("CreateUser_Fail", "ip", loadRealIp(ctx), "userId", session.UserId, "cause", er.Error())
		return nil, transport.WrapInteralServerError(er)
	}

	if err := dTx.Commit(); err != nil {
		foree_logger.Logger.Error("CreateUser_Fail", "ip", loadRealIp(ctx), "userId", session.UserId, "cause", err.Error())
		return nil, transport.WrapInteralServerError(err)
	}

	foree_logger.Logger.Info("CreateUser_Success", "ip", loadRealIp(ctx), "userId", session.UserId)

	userGroup, er := a.userGroupRepo.GetUniqueUserGroupByOwnerId(newUser.ID)
	if er != nil {
		foree_logger.Logger.Error("CreateUser_Fail", "ip", loadRealIp(ctx), "userId", session.UserId, "cause", er.Error())
		return nil, transport.WrapInteralServerError(er)
	}

	// Get Permission.
	rolePermissions, pErr := a.rolePermissionRepo.GetAllEnabledRolePermissionByRoleName(userGroup.RoleGroup)
	if pErr != nil {
		foree_logger.Logger.Error("CreateUser_Fail", "ip", loadRealIp(ctx), "userId", session.UserId, "cause", pErr.Error())
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
		foree_logger.Logger.Error("CreateUser_Fail", "ip", loadRealIp(ctx), "userId", session.UserId, "cause", sessionErr.Error())
		return nil, transport.WrapInteralServerError(sessionErr)
	}

	go func() {
		// Create default Interac Account for the user.
		now := time.Now()
		acc := account.InteracAccount{
			FirstName:        newUser.FirstName,
			MiddleName:       newUser.MiddleName,
			LastName:         newUser.LastName,
			Address1:         newUser.Address1,
			Address2:         newUser.Address2,
			City:             newUser.City,
			Province:         newUser.Province,
			Country:          newUser.Country,
			PostalCode:       newUser.PostalCode,
			PhoneNumber:      newUser.PhoneNumber,
			Email:            newUser.Email,
			OwnerId:          newUser.ID,
			Status:           account.AccountStatusActive,
			LatestActivityAt: &now,
		}
		_, derr := a.interacAccountRepo.InsertInteracAccount(context.TODO(), acc)
		if derr != nil {
			foree_logger.Logger.Error("Default_Interac_Account_Fail", "ip", loadRealIp(ctx), "userId", session.UserId, "cause", derr.Error())
		}
	}()

	go a.rewardOnboard(newUser)

	return NewUserDTO(updateSession), nil
}

// TODO: Login protection on peak volume.
func (a *AuthService) Login(ctx context.Context, req LoginReq) (*UserDTO, transport.HError) {
	// Delete previous token if exists.
	a.sessionRepo.Delete(req.SessionId)

	// Verify email and password
	ep, err := a.emailPasswordRepo.GetUniqueEmailPasswdByEmail(req.Email)

	if err != nil {
		foree_logger.Logger.Error("Login_Fail", "ip", loadRealIp(ctx), "email", req.Email, "cause", err.Error())
		return nil, transport.WrapInteralServerError(err)
	}

	if ep == nil {
		foree_logger.Logger.Warn("Login_Fail", "ip", loadRealIp(ctx), "email", req.Email, "cause", fmt.Sprintf("email `%s` not found", req.Email))
		return nil, transport.NewFormError("Invaild signup", "email", "invalid email")
	}

	if ep.Status == auth.EPStatusDelete {
		foree_logger.Logger.Warn("Login_Fail", "ip", loadRealIp(ctx), "email", req.Email, "cause", fmt.Sprintf("email `%s` is in `%s` status", ep.Email, ep.Status))
		return nil, transport.NewFormError("Invalid login request", "email", "invalid email")
	}

	if ep.Status == auth.EPStatusSuspend {
		foree_logger.Logger.Warn("Login_Fail", "ip", loadRealIp(ctx), "email", req.Email, "cause", fmt.Sprintf("email `%s` is in `%s` status", ep.Email, ep.Status))
		return nil, transport.NewFormError("Invalid login request", "email", "your account is suspend. please contact us.")
	}

	if ep.LoginAttempts > maxLoginAttempts {
		foree_logger.Logger.Warn("Login_Fail", "ip", loadRealIp(ctx), "email", req.Email, "cause", fmt.Sprintf("email `%s` has try `%v` times", ep.Email, ep.LoginAttempts))
		return nil, transport.NewFormError("Invalid login request", "password", "max login attempts reached. please contact us.")
	}

	ok := auth.ComparePasswords(req.Password, []byte(ep.Passwd))
	if !ok {
		foree_logger.Logger.Warn("Login_Fail", "ip", loadRealIp(ctx), "email", req.Email, "cause", "invalid password")
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
		foree_logger.Logger.Error("Login_Fail", "ip", loadRealIp(ctx), "email", req.Email, "cause", err.Error())
		return nil, transport.WrapInteralServerError(err)
	}
	// User must exists
	if user == nil {
		foree_logger.Logger.Error("Login_Fail", "ip", loadRealIp(ctx), "email", req.Email, "cause", fmt.Sprintf("owner `%v` no found", ep.OwnerId))
		return nil, transport.NewInteralServerError("User `%v` do not exists", ep.OwnerId)
	}

	userGroup, er := a.userGroupRepo.GetUniqueUserGroupByOwnerId(user.ID)
	if er != nil {
		foree_logger.Logger.Error("Login_Fail", "ip", loadRealIp(ctx), "email", req.Email, "user", user.ID, "cause", er.Error())
		return nil, transport.WrapInteralServerError(er)
	}
	//User group must exists
	if userGroup == nil {
		foree_logger.Logger.Error("Login_Fail", "ip", loadRealIp(ctx), "email", req.Email, "cause", fmt.Sprintf("userGroup not found with owener `%v`", ep.OwnerId))
		return nil, transport.NewInteralServerError("User `%v` do not exists", ep.OwnerId)
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
		foree_logger.Logger.Error("Login_Fail", "ip", loadRealIp(ctx), "email", req.Email, "cause", "fail to insert session to session repo")
		return nil, transport.WrapInteralServerError(err)
	}
	session := a.sessionRepo.GetSessionUniqueById(sessionId)
	if session == nil {
		foree_logger.Logger.Error("Login_Fail", "ip", loadRealIp(ctx), "email", req.Email, "cause", "fail to get session from session repo")
		return nil, transport.NewInteralServerError("sesson `%s` not found", sessionId)
	}

	foree_logger.Logger.Info("Login_Success", "ip", loadRealIp(ctx), "email", req.Email, "userAgent", loadUserAgent(ctx))
	return NewUserDTO(session), nil
}

func (a *AuthService) ForgetPassword(ctx context.Context, req ForgetPasswordReq) (any, transport.HError) {
	// Verify email and password
	ep, err := a.emailPasswordRepo.GetUniqueEmailPasswdByEmail(req.Email)

	if err != nil {
		foree_logger.Logger.Error("ForgetPassword_Fail", "ip", loadRealIp(ctx), "email", req.Email, "cause", err.Error())
		return nil, transport.WrapInteralServerError(err)

	}

	if ep == nil {
		foree_logger.Logger.Warn("ForgetPassword_Fail", "ip", loadRealIp(ctx), "email", req.Email, "cause", fmt.Sprintf("email `%s` not found", req.Email))
		return nil, transport.NewFormError("Invaild forget password request", "email", "invalid email")
	}

	newEP := *ep
	retrieveTokenExpiredAt := time.Now().Add(5 * time.Minute)
	newEP.RetrieveToken = auth.GenerateVerifyCode()
	newEP.RetrieveTokenExpiredAt = &retrieveTokenExpiredAt

	err = a.emailPasswordRepo.UpdateEmailPasswdByEmail(context.Background(), newEP)
	if err != nil {
		foree_logger.Logger.Error("ForgetPassword_Fail", "ip", loadRealIp(ctx), "email", req.Email, "cause", err.Error())
		transport.WrapInteralServerError(err)
	}

	foree_logger.Logger.Info("ForgetPassword_Success", "ip", loadRealIp(ctx), "email", req.Email)

	//TODO: send email with code

	return nil, nil
}

func (a *AuthService) ForgetPasswordUpdate(ctx context.Context, req ForgetPasswordUpdateReq) (any, transport.HError) {
	ep, err := a.emailPasswordRepo.GetUniqueEmailPasswdByEmail(req.Email)

	if err != nil {
		foree_logger.Logger.Error("ForgetPasswordUpdate_Fail", "ip", loadRealIp(ctx), "email", req.Email, "cause", err.Error())
		return nil, transport.WrapInteralServerError(err)

	}

	if ep == nil {
		foree_logger.Logger.Warn("ForgetPasswordUpdate_Fail", "ip", loadRealIp(ctx), "email", req.Email, "cause", fmt.Sprintf("email `%s` not found", req.Email))
		return nil, transport.NewFormError("Invaild forget password", "email", "invalid email")
	}

	if ep.RetrieveTokenExpiredAt == nil && ep.RetrieveTokenExpiredAt.Before(time.Now()) {
		foree_logger.Logger.Warn("ForgetPasswordUpdate_Fail", "ip", loadRealIp(ctx), "email", req.Email, "cause", "retrieve code expire")
		return nil, transport.NewFormError("Invaild forget password", "retrieveCode", "retrieve code expired")
	}

	if ep.RetrieveToken != req.RetrieveCode {
		foree_logger.Logger.Warn("ForgetPasswordUpdate_Fail", "ip", loadRealIp(ctx), "email", req.Email, "cause", "invalid retrieve code")
		return nil, transport.NewFormError("Invaild forget password", "retrieveCode", "invalid retrieve code")
	}

	hashedPasswd, err := auth.HashPassword(req.NewPassword)
	if err != nil {
		foree_logger.Logger.Error("ForgetPasswordUpdate_Fail", "ip", loadRealIp(ctx), "email", req.Email, "cause", err.Error())
		return nil, transport.WrapInteralServerError(err)
	}

	newEP := *ep
	newEP.Passwd = hashedPasswd

	err = a.emailPasswordRepo.UpdateEmailPasswdByEmail(context.Background(), newEP)

	if err != nil {
		foree_logger.Logger.Error("ForgetPasswordUpdate_Fail", "ip", loadRealIp(ctx), "email", req.Email, "cause", err.Error())
		transport.WrapInteralServerError(err)
	}

	//TODO: send email to user for notice

	foree_logger.Logger.Info("ForgetPasswordUpdate_Success", "ip", loadRealIp(ctx), "email", req.Email)
	return nil, nil
}

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

	return NewUserDTO(session), nil
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
		return session, err
	}
	for _, p := range session.RolePermissions {
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
	if req.ReferralCode == "" {
		return
	}

	referral, err := a.referralRepo.GetUniqueReferralByReferralCode(req.ReferralCode)
	if err != nil {
		foree_logger.Logger.Error("Link_Referer_Fail", "userId", registerUser.ID, "ReferralCode", req.ReferralCode, "cause", err.Error())
		return
	}
	if referral == nil {
		foree_logger.Logger.Warn("Link_Referer_Fail", "userId", registerUser.ID, "ReferralCode", req.ReferralCode, "cause", "unknown ReferralCode")
		return
	}
	if referral.RefereeId != 0 {
		foree_logger.Logger.Warn("Link_Referer_Fail", "userId", registerUser.ID, "ReferralCode", req.ReferralCode, "cause", "unknown ReferralCode")
		return
	}

	newReferral := *referral
	newReferral.RefereeId = registerUser.ID
	now := time.Now()
	newReferral.AcceptAt = &now

	err = a.referralRepo.UpdateReferralByReferralCode(newReferral)
	if err != nil {
		foree_logger.Logger.Error("Link_Referer_Fail", "userId", registerUser.ID, "ReferralCode", req.ReferralCode, "cause", err.Error())
		return
	}
	foree_logger.Logger.Info("Link_Referer_success", "userId", registerUser.ID, "ReferrerId", referral.ReferrerId)
}

// TODO: configure to turn targo.
func (a *AuthService) rewardReferer(registerUser auth.User) {
	referral, _ := a.referralRepo.GetUniqueReferralByRefereeId(registerUser.ID)
	if referral == nil {
		return
	}

	promotion, _ := a.getPromotion(PromotionReferral, promotionCacheTimeout)

	if promotion == nil || !promotion.IsValid() {
		return
	}

	reward := transaction.Reward{
		Type:        transaction.RewardTypeReferal,
		Description: fmt.Sprintf("Referral reward by %v %v", registerUser.FirstName, registerUser.LastName),
		Amt:         promotion.Amt,
		OwnerId:     referral.ReferrerId,
		ExpireAt:    time.Now().Add(time.Hour * 24 * 180),
	}

	_, err := a.rewardRepo.InsertReward(context.TODO(), reward)
	if err != nil {
		foree_logger.Logger.Error("Referral_Reward_Fail", "refereeId", registerUser.ID, "referrerId", referral.ReferrerId, "cause", err.Error())
	}
}

func (a *AuthService) rewardOnboard(registerUser auth.User) {

	gift, _ := a.getPromotion(PromotionOnboard, promotionCacheTimeout)

	if gift == nil || !gift.IsValid() {
		return
	}

	reward := transaction.Reward{
		Type:        transaction.RewardTypeReferal,
		Description: "Onboard reward",
		Amt:         gift.Amt,
		OwnerId:     registerUser.ID,
		ExpireAt:    time.Now().Add(time.Hour * 24 * 180),
	}

	_, err := a.rewardRepo.InsertReward(context.TODO(), reward)
	if err != nil {
		foree_logger.Logger.Error("Onboard_Reward_Fail", "userId", registerUser.ID, "cause", err.Error())
	}
}

// TODO: using atomic interger to limit peak volumn
func (a *AuthService) getPromotion(promotionCode string, validIn time.Duration) (*promotion.Promotion, error) {
	a.promotionCacheRWLock.RLock()
	promotionCache, ok := a.promotionCache[promotionCode]
	a.promotionCacheRWLock.RUnlock()

	if ok && promotionCache.createdAt.Add(validIn).After(time.Now()) {
		return &promotionCache.item, nil
	}

	gift, err := a.promotionRepo.GetUniquePromotionByCode(context.TODO(), promotionCode)
	if err != nil {
		foree_logger.Logger.Error("Promotion_Fail", "promotionCode", promotionCode, "cause", err.Error())
		return nil, err
	}

	if gift != nil {
		foree_logger.Logger.Warn("Promotion_Fail", "promotionCode", promotionCode, "cause", "gift no found")
		return nil, fmt.Errorf("Promotion no found with code `%v`", promotionCode)
	}

	// Update gift
	// Make sure at least one thread can update the cache.
	func() {
		a.promotionCacheUpdateLock.TryLock()
		defer a.promotionCacheUpdateLock.Unlock()
		a.promotionCacheRWLock.Lock()
		defer a.promotionCacheRWLock.Unlock()
		a.promotionCache[promotionCode] = CacheItem[promotion.Promotion]{
			item:      *gift,
			createdAt: time.Now(),
		}
	}()

	return gift, nil
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
