package foree_service

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"xue.io/go-pay/app/foree/account"
	foree_auth "xue.io/go-pay/app/foree/auth"
	foree_constant "xue.io/go-pay/app/foree/constant"
	"xue.io/go-pay/app/foree/logger"
	"xue.io/go-pay/app/foree/promotion"
	"xue.io/go-pay/app/foree/referral"
	"xue.io/go-pay/app/foree/transaction"
	"xue.io/go-pay/auth"
	"xue.io/go-pay/constant"
	"xue.io/go-pay/server/transport"
)

const maxLoginAttempts = 4
const OnboardGift = "ONBOARD_GIFT"
const ReferralGift = "REFERRAL_GIFT"
const giftCacheTimeout = 15 * time.Minute

func NewAuthService(
	db *sql.DB,
	sessionRepo *auth.SessionRepo,
	userRepo *auth.UserRepo,
	emailPasswordRepo *auth.EmailPasswdRepo,
	rolePermissionRepo *auth.RolePermissionRepo,
	userIdentificationRepo *foree_auth.UserIdentificationRepo,
	interacAccountRepo *account.InteracAccountRepo,
	userGroupRepo *auth.UserGroupRepo,
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
	referralRepo           *referral.ReferralRepo
	rewardRepo             *transaction.RewardRepo
	giftRepo               *promotion.GiftRepo
	giftCache              map[string]CacheItem[promotion.Gift]
	giftCacheRWLock        sync.RWMutex
	giftCacheUpdateLock    sync.RWMutex
}

func (a *AuthService) SignUp(ctx context.Context, req SignUpReq) (*UserDTO, transport.HError) {
	// Check if email already exists.
	oldEmail, err := a.emailPasswordRepo.GetUniqueEmailPasswdByEmail(req.Email)
	if err != nil {
		logger.Logger.Error("Signup_Fail", "ip", loadRealIp(ctx), "email", req.Email, "cause", err.Error())
		return nil, transport.WrapInteralServerError(err)
	}

	if oldEmail != nil {
		logger.Logger.Warn("Signup_Fail", "ip", loadRealIp(ctx), "email", req.Email, "cause", "email already exists")
		return nil, transport.NewFormError("Invaild signup", "email", "Duplicate email")
	}

	// Hashing password.
	hashedPasswd, err := auth.HashPassword(req.Password)
	if err != nil {
		logger.Logger.Error("Signup_Fail", "ip", loadRealIp(ctx), "email", req.Email, "cause", err.Error())
		return nil, transport.WrapInteralServerError(err)
	}

	// Start DB transaction
	dTx, err := a.db.Begin()
	if err != nil {
		logger.Logger.Error("Signup_Fail", "ip", loadRealIp(ctx), "email", req.Email, "cause", err.Error())
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
		logger.Logger.Error("Signup_Fail", "ip", loadRealIp(ctx), "email", req.Email, "cause", err.Error())
		return nil, transport.WrapInteralServerError(err)
	}

	user, err := a.userRepo.GetUniqueUserById(userId)
	if err != nil {
		dTx.Rollback()
		logger.Logger.Error("Signup_Fail", "ip", loadRealIp(ctx), "email", req.Email, "cause", err.Error())
		return nil, transport.WrapInteralServerError(err)
	}

	if user == nil {
		dTx.Rollback()
		logger.Logger.Error("Signup_Fail", "ip", loadRealIp(ctx), "email", req.Email, "userId", userId, "cause", "unable to get user")
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
		logger.Logger.Error("Signup_Fail", "ip", loadRealIp(ctx), "email", req.Email, "cause", err.Error())
		return nil, transport.WrapInteralServerError(err)
	}

	ep, err := a.emailPasswordRepo.GetUniqueEmailPasswdById(id)

	if err != nil {
		dTx.Rollback()
		logger.Logger.Error("Signup_Fail", "ip", loadRealIp(ctx), "email", req.Email, "cause", err.Error())
		return nil, transport.WrapInteralServerError(err)
	}

	if ep == nil {
		dTx.Rollback()
		logger.Logger.Error("Signup_Fail", "ip", loadRealIp(ctx), "email", req.Email, "emailPasswordId", id, "cause", "unable to get created emailPassword")
		return nil, transport.NewInteralServerError("unable to get EmailPasswd with id: `%v`", id)
	}

	if err = dTx.Commit(); err != nil {
		logger.Logger.Error("Signup_Fail", "ip", loadRealIp(ctx), "email", req.Email, "cause", err.Error())
		return nil, transport.WrapInteralServerError(err)
	}

	go a.linkReferer(*user, req)

	//TODO: send email. by goroutine

	sessionId, err := a.sessionRepo.InsertSession(auth.Session{
		UserId:      user.ID,
		User:        user,
		EmailPasswd: ep,
	})
	if err != nil {
		logger.Logger.Error("Signup_Fail", "ip", loadRealIp(ctx), "email", req.Email, "cause", err.Error())
		return nil, transport.WrapInteralServerError(err)
	}

	session := a.sessionRepo.GetSessionUniqueById(sessionId)
	if session == nil {
		logger.Logger.Error("Signup_Fail", "ip", loadRealIp(ctx), "email", req.Email, "cause", "unable to get created session")
		return nil, transport.NewInteralServerError("sesson `%s` not found", sessionId)
	}
	return NewUserDTO(session), nil
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
	_, derr := a.interacAccountRepo.InsertInteracAccount(ctx, acc)
	if derr != nil {
		return nil, transport.WrapInteralServerError(derr)
	}
	return updateSession, nil
}

// TODO: Login protection on peak volume.
func (a *AuthService) Login(ctx context.Context, req LoginReq) (*UserDTO, transport.HError) {
	// Delete previous token if exists.
	a.sessionRepo.Delete(req.SessionId)

	// Verify email and password
	ep, err := a.emailPasswordRepo.GetUniqueEmailPasswdByEmail(req.Email)

	if err != nil {
		logger.Logger.Error("Login_Fail", "ip", loadRealIp(ctx), "email", req.Email, "cause", err.Error())
		return nil, transport.WrapInteralServerError(err)
	}

	if ep == nil {
		logger.Logger.Warn("Login_Fail", "ip", loadRealIp(ctx), "email", req.Email, "cause", fmt.Sprintf("email `%s` not found", req.Email))
		return nil, transport.NewFormError("Invaild signup", "email", "invalid email")
	}

	if ep.Status == auth.EPStatusDelete {
		logger.Logger.Warn("Login_Fail", "ip", loadRealIp(ctx), "email", req.Email, "cause", fmt.Sprintf("email `%s` is in `%s` status", ep.Email, ep.Status))
		return nil, transport.NewFormError("Invalid login request", "email", "invalid email")
	}

	if ep.Status == auth.EPStatusSuspend {
		logger.Logger.Warn("Login_Fail", "ip", loadRealIp(ctx), "email", req.Email, "cause", fmt.Sprintf("email `%s` is in `%s` status", ep.Email, ep.Status))
		return nil, transport.NewFormError("Invalid login request", "email", "your account is suspend. please contact us.")
	}

	if ep.LoginAttempts > maxLoginAttempts {
		logger.Logger.Warn("Login_Fail", "ip", loadRealIp(ctx), "email", req.Email, "cause", fmt.Sprintf("email `%s` has try `%v` times", ep.Email, ep.LoginAttempts))
		return nil, transport.NewFormError("Invalid login request", "password", "max login attempts reached. please contact us.")
	}

	ok := auth.ComparePasswords(req.Password, []byte(ep.Passwd))
	if !ok {
		logger.Logger.Warn("Login_Fail", "ip", loadRealIp(ctx), "email", req.Email, "cause", "invalid password")
		go func() {
			newEP := *ep
			newEP.LoginAttempts += 1
			if err := a.emailPasswordRepo.UpdateEmailPasswdByEmail(ctx, newEP); err != nil {
				logger.Logger.Error("Login_Attempts_Update", "email", req.Email, "cause", err.Error())
			}
		}()
		return nil, transport.NewFormError("Invaild signup", "password", "Invalid password")
	}

	// Load user(user must exist, but not necessary to be active)
	user, err := a.userRepo.GetUniqueUserById(ep.OwnerId)
	if err != nil {
		logger.Logger.Error("Login_Fail", "ip", loadRealIp(ctx), "email", req.Email, "cause", err.Error())
		return nil, transport.WrapInteralServerError(err)
	}
	// User must exists
	if user == nil {
		logger.Logger.Error("Login_Fail", "ip", loadRealIp(ctx), "email", req.Email, "cause", fmt.Sprintf("owner `%v` no found", ep.OwnerId))
		return nil, transport.NewInteralServerError("User `%v` do not exists", ep.OwnerId)
	}

	userGroup, er := a.userGroupRepo.GetUniqueUserGroupByOwnerId(user.ID)
	if er != nil {
		logger.Logger.Error("Login_Fail", "ip", loadRealIp(ctx), "email", req.Email, "user", user.ID, "cause", er.Error())
		return nil, transport.WrapInteralServerError(er)
	}
	//User group must exists
	if userGroup == nil {
		logger.Logger.Error("Login_Fail", "ip", loadRealIp(ctx), "email", req.Email, "cause", fmt.Sprintf("userGroup not found with owener `%v`", ep.OwnerId))
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
		logger.Logger.Error("Login_Fail", "ip", loadRealIp(ctx), "email", req.Email, "cause", "fail to insert session to session repo")
		return nil, transport.WrapInteralServerError(err)
	}
	session := a.sessionRepo.GetSessionUniqueById(sessionId)
	if session == nil {
		logger.Logger.Error("Login_Fail", "ip", loadRealIp(ctx), "email", req.Email, "cause", "fail to get session from session repo")
		return nil, transport.NewInteralServerError("sesson `%s` not found", sessionId)
	}

	logger.Logger.Info("Login_Success", "ip", loadRealIp(ctx), "email", req.Email, "userAgent", loadUserAgent(ctx))
	return NewUserDTO(session), nil
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

func (a *AuthService) linkReferer(registerUser auth.User, req SignUpReq) {
	if req.ReferralCode == "" {
		return
	}

	referral, err := a.referralRepo.GetUniqueReferralByReferralCode(req.ReferralCode)
	if err != nil {
		logger.Logger.Error("Link_Referer_Fail", "userId", registerUser.ID, "ReferralCode", req.ReferralCode, "cause", err.Error())
		return
	}
	if referral == nil {
		logger.Logger.Warn("Link_Referer_Fail", "userId", registerUser.ID, "ReferralCode", req.ReferralCode, "cause", "unknown ReferralCode")
		return
	}
	if referral.RefereeId != 0 {
		logger.Logger.Warn("Link_Referer_Fail", "userId", registerUser.ID, "ReferralCode", req.ReferralCode, "cause", "unknown ReferralCode")
		return
	}

	newReferral := *referral
	newReferral.RefereeId = registerUser.ID
	newReferral.AcceptAt = time.Now()

	err = a.referralRepo.UpdateReferralByReferralCode(newReferral)
	if err != nil {
		logger.Logger.Error("Link_Referer_Fail", "userId", registerUser.ID, "ReferralCode", req.ReferralCode, "cause", err.Error())
		return
	}
	logger.Logger.Info("Link_Referer_success", "userId", registerUser.ID, "ReferrerId", referral.ReferrerId)
}

// TODO: configure to turn targo.
func (a *AuthService) rewardReferer(registerUser auth.User) {
	referral, _ := a.referralRepo.GetUniqueReferralByRefereeId(registerUser.ID)
	if referral == nil {
		return
	}

	gift, _ := a.getGift(ReferralGift, giftCacheTimeout)

	if gift == nil || !gift.IsValid() {
		return
	}

	reward := transaction.Reward{
		Type:        transaction.RewardTypeReferal,
		Description: fmt.Sprintf("Referral reward by %v %v", registerUser.FirstName, registerUser.LastName),
		Amt:         gift.Amt,
		OwnerId:     referral.ReferrerId,
		ExpireAt:    time.Now().Add(time.Hour * 24 * 180),
	}

	_, err := a.rewardRepo.InsertReward(context.TODO(), reward)
	if err != nil {
		logger.Logger.Error("Referral_Reward_Fail", "refereeId", registerUser.ID, "referrerId", referral.ReferrerId, "cause", err.Error())
	}
}

func (a *AuthService) rewardOnboard(registerUser auth.User) {

	gift, _ := a.getGift(OnboardGift, giftCacheTimeout)

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
		logger.Logger.Error("Onboard_Reward_Fail", "userId", registerUser.ID, "cause", err.Error())
	}
}

// TODO: using atomic interger to limit peak volumn
func (a *AuthService) getGift(giftCode string, validIn time.Duration) (*promotion.Gift, error) {
	a.giftCacheRWLock.RLock()
	giftCache, ok := a.giftCache[giftCode]
	a.giftCacheRWLock.RUnlock()

	if ok && giftCache.createdAt.Add(validIn).After(time.Now()) {
		return &giftCache.item, nil
	}

	gift, err := a.giftRepo.GetUniqueGiftByCode(context.TODO(), giftCode)
	if err != nil {
		logger.Logger.Error("Gift_Fail", "giftCode", giftCode, "cause", err.Error())
		return nil, err
	}

	if gift != nil {
		logger.Logger.Warn("Gift_Fail", "giftCode", giftCode, "cause", "gift no found")
		return nil, fmt.Errorf("Gift no found with giftCode `%v`", giftCode)
	}

	// Update gift
	// Make sure at least one thread can update the cache.
	func() {
		a.giftCacheUpdateLock.TryLock()
		defer a.giftCacheUpdateLock.Unlock()
		a.giftCacheRWLock.Lock()
		defer a.giftCacheRWLock.Unlock()
		a.giftCache[giftCode] = CacheItem[promotion.Gift]{
			item:      *gift,
			createdAt: time.Now(),
		}
	}()

	return gift, nil
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

func loadXForwardFor(ctx context.Context) string {
	req, ok := ctx.Value(constant.CKHttpRequest).(*http.Request)
	if !ok {
		return ""
	}
	return req.Header.Get("X-Forwarded-For")
}

func loadIp(ctx context.Context) string {
	req, ok := ctx.Value(constant.CKHttpRequest).(*http.Request)
	if !ok {
		return ""
	}
	return req.RemoteAddr
}

func loadRealIp(ctx context.Context) string {
	xforward := loadXForwardFor(ctx)
	if xforward != "" && len(strings.Split(xforward, ",")) == 0 {
		return loadIp(ctx)
	}
	return strings.Split(xforward, ",")[0]
}

func loadUserAgent(ctx context.Context) string {
	req, ok := ctx.Value(constant.CKHttpRequest).(*http.Request)
	if !ok {
		return ""
	}
	return req.Header.Get("User-Agent")
}
