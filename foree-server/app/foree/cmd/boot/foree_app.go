package foree_boot

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"xue.io/go-pay/app/foree/account"
	foree_router "xue.io/go-pay/app/foree/app_router"
	foree_service "xue.io/go-pay/app/foree/app_service"
	foree_auth "xue.io/go-pay/app/foree/auth"
	foree_config "xue.io/go-pay/app/foree/cmd/config"
	foree_logger "xue.io/go-pay/app/foree/logger"
	"xue.io/go-pay/app/foree/promotion"
	"xue.io/go-pay/app/foree/referral"
	"xue.io/go-pay/app/foree/sys_router"
	"xue.io/go-pay/app/foree/transaction"
	"xue.io/go-pay/auth"
	"xue.io/go-pay/config"
	ms "xue.io/go-pay/db/mysql"
	"xue.io/go-pay/logger"
	"xue.io/go-pay/partner/idm"
	"xue.io/go-pay/partner/nbp"
	"xue.io/go-pay/partner/scotia"
)

type ForeeApp struct {
	envFilePath              string
	db                       *sql.DB
	userRepo                 *auth.UserRepo
	userGroupRepo            *auth.UserGroupRepo
	sessionRepo              *auth.SessionRepo
	userSettingRepo          *auth.UserSettingRepo
	emailPasswdRepo          *auth.EmailPasswdRepo
	rolePermissionRepo       *auth.RolePermissionRepo
	contactAccountRepo       *account.ContactAccountRepo
	interacAccountRepo       *account.InteracAccountRepo
	userExtraRepo            *foree_auth.UserExtraRepo
	userIdnetificationRepo   *foree_auth.UserIdentificationRepo
	referralRepo             *referral.ReferralRepo
	dailyTxLimitRepo         *transaction.DailyTxLimitRepo
	txLimitRepo              *transaction.TxLimitRepo
	feeRepo                  *transaction.FeeRepo
	feeJointRepo             *transaction.FeeJointRepo
	rewardRepo               *promotion.RewardRepo
	rateRepo                 *transaction.RateRepo
	idmTxRepo                *transaction.IdmTxRepo
	idmRepo                  *transaction.IDMComplianceRepo
	foreeTxRepo              *transaction.ForeeTxRepo
	txHistoryRepo            *transaction.TxHistoryRepo
	interacCITxRepo          *transaction.InteracCITxRepo
	interacRefundTxRepo      *transaction.ForeeRefundTxRepo
	nbpCOTxRepo              *transaction.NBPCOTxRepo
	txQuoteRepo              *transaction.TxQuoteRepo
	txSummaryRepo            *transaction.TxSummaryRepo
	promotionRepo            *promotion.PromotionRepo
	authService              *foree_service.AuthService
	accountService           *foree_service.AccountService
	feeService               *foree_service.FeeService
	rateService              *foree_service.RateService
	txLimitService           *foree_service.TxLimitService
	transactionService       *foree_service.TransactionService
	promotionService         *foree_service.PromotionService
	promotionRewardJointRepo *promotion.PromotionRewardJointRepo
	scotiaClient             scotia.ScotiaClient
	idmClient                idm.IDMClient
	nbpClient                nbp.NBPClient
	txProcessor              *foree_service.TxProcessor
	interacTxProcessor       *foree_service.InteracTxProcessor
	idmTxProcessor           *foree_service.IDMTxProcessor
	nbpTxProcessor           *foree_service.NBPTxProcessor
	accountRouter            *foree_router.AccountRouter
	authRouter               *foree_router.AuthRouter
	transactionRouter        *foree_router.TransactionRouter
	sysRouter                *sys_router.SystemRouter
	generateRouter           *foree_router.GeneralRouter
}

func (app *ForeeApp) Boot(envFilePath string) error {
	app.envFilePath = envFilePath
	var cfg foree_config.ForeeLocalConfig
	if err := config.LoadFromFile(&cfg, envFilePath); err != nil {
		return err
	}

	//Initial Logger
	l, err := logger.NewZapLogger("debug", "/tmp/foree/foree.log")
	if err != nil {
		panic(err)
	}
	foree_logger.Logger = l

	//Initial DB
	db, err := ms.NewMysqlPool(mysql.Config{
		Addr:                 fmt.Sprintf("%s:%s", cfg.MysqlDBHost, cfg.MysqlDBPort),
		DBName:               cfg.MysqlDBName,
		User:                 cfg.MysqlDBUser,
		Passwd:               cfg.MysqlDBPasswd,
		AllowNativePasswords: true,
		ParseTime:            true,
	}, 1, 1)

	if err != nil {
		return err
	}
	app.db = db

	// Initial Repos
	app.userRepo = auth.NewUserRepo(db)
	app.userGroupRepo = auth.NewUserGroupRepo(db)
	app.userSettingRepo = auth.NewUserSettingRepo(db)
	app.sessionRepo = auth.NewDefaultSessionRepo(db)
	app.emailPasswdRepo = auth.NewEmailPasswdRepo(db)
	app.rolePermissionRepo = auth.NewRolePermission(db)
	app.contactAccountRepo = account.NewContactAccountRepo(db)
	app.interacAccountRepo = account.NewInteracAccountRepo(db)
	app.userExtraRepo = foree_auth.NewUserExtraRepo(db)
	app.userIdnetificationRepo = foree_auth.NewUserIdentificationRepo(db)
	app.referralRepo = referral.NewReferralRepo(db)
	app.dailyTxLimitRepo = transaction.NewDailyTxLimitRepo(db)
	app.feeRepo = transaction.NewFeeRepo(db)
	app.feeJointRepo = transaction.NewFeeJointRepo(db)
	app.rateRepo = transaction.NewRateRepo(db)
	app.rewardRepo = promotion.NewRewardRepo(db)
	app.idmTxRepo = transaction.NewIdmTxRepo(db)
	app.idmRepo = transaction.NewIDMComplianceRepo(db)
	app.foreeTxRepo = transaction.NewForeeTxRepo(db)
	app.txHistoryRepo = transaction.NewTxHistoryRepo(db)
	app.interacCITxRepo = transaction.NewInteracCITxRepo(db)
	app.interacRefundTxRepo = transaction.NewForeeRefundTxRepo(db)
	app.nbpCOTxRepo = transaction.NewNBPCOTxRepo(db)
	app.txQuoteRepo = transaction.NewTxQuoteRepo(5, 2048)
	app.txSummaryRepo = transaction.NewTxSummaryRepo(db)
	app.promotionRepo = promotion.NewPromotionRepo(db)
	app.promotionRewardJointRepo = promotion.NewPromotionRewardJointRepo(db)
	app.txLimitRepo = transaction.NewTxLimitRepo(db)

	//Initial vendors
	app.scotiaClient = scotia.NewMockScotiaClient()
	app.idmClient = idm.NewMockIDMClient()
	app.nbpClient = nbp.NewMockNBPClient()

	//Initial transaction processors
	app.txProcessor = foree_service.NewTxProcessor(
		app.db,
		app.interacCITxRepo,
		app.nbpCOTxRepo,
		app.idmTxRepo,
		app.txHistoryRepo,
		app.txSummaryRepo,
		app.foreeTxRepo,
		app.interacRefundTxRepo,
		app.rewardRepo,
		app.dailyTxLimitRepo,
		app.userRepo,
		app.contactAccountRepo,
		app.interacAccountRepo,
	)
	app.interacTxProcessor = foree_service.NewInteracTxProcessor(
		app.db,
		foree_service.ScotiaProfile{},
		app.scotiaClient,
		app.interacCITxRepo,
		app.foreeTxRepo,
		app.txSummaryRepo,
		app.txProcessor,
	)
	app.idmTxProcessor = foree_service.NewIDMTxProcessor(
		app.db,
		app.foreeTxRepo,
		app.idmTxRepo,
		app.idmClient,
		app.txProcessor,
	)
	app.nbpTxProcessor = foree_service.NewNBPTxProcessor(
		app.db,
		app.foreeTxRepo,
		app.txProcessor,
		app.nbpCOTxRepo,
		app.nbpClient,
		app.userExtraRepo,
		app.userIdnetificationRepo,
	)

	app.txProcessor.SetInteracTxProcessor(app.interacTxProcessor)
	app.txProcessor.SetIDMTxProcessor(app.idmTxProcessor)
	app.txProcessor.SetNBPTxProcessor(app.nbpTxProcessor)

	app.promotionService = foree_service.NewPromotionService(app.db, app.promotionRepo, app.rewardRepo, app.referralRepo, app.promotionRewardJointRepo)

	//Initial service
	app.authService = foree_service.NewAuthService(
		db, app.sessionRepo,
		app.userRepo,
		app.emailPasswdRepo,
		app.rolePermissionRepo,
		app.userIdnetificationRepo,
		app.interacAccountRepo,
		app.userGroupRepo,
		app.userExtraRepo,
		app.userSettingRepo,
		app.referralRepo,
		app.promotionService,
	)

	app.accountService = foree_service.NewAccountService(
		app.authService,
		app.contactAccountRepo,
		app.interacAccountRepo,
	)

	app.rateService = foree_service.NewRateService(app.rateRepo)
	app.feeService = foree_service.NewFeeService(app.feeRepo)
	app.txLimitService = foree_service.NewTxLimitService(app.txLimitRepo, app.dailyTxLimitRepo)

	app.transactionService = foree_service.NewTransactionService(
		db,
		app.authService,
		app.userGroupRepo,
		app.foreeTxRepo,
		app.txSummaryRepo,
		app.txQuoteRepo,
		app.rewardRepo,
		app.contactAccountRepo,
		app.interacAccountRepo,
		app.feeJointRepo,
		app.rateService,
		app.feeService,
		app.txLimitService,
		app.txProcessor,
		app.scotiaClient,
		app.nbpClient,
	)

	//Initial router
	app.accountRouter = foree_router.NewAccountRouter(app.authService, app.accountService)
	app.authRouter = foree_router.NewAuthRouter(app.authService)
	app.transactionRouter = foree_router.NewTransactionRouter(app.authService, app.transactionService)
	app.generateRouter = foree_router.NewGeneralRouter()

	router := mux.NewRouter()
	subrouter := router.PathPrefix("/app/v1").Subrouter()

	app.accountRouter.RegisterRouter(subrouter)
	app.authRouter.RegisterRouter(subrouter)
	app.transactionRouter.RegisterRouter(subrouter)
	app.generateRouter.RegisterRouter(subrouter)

	sysSubrouter := router.PathPrefix("/sys/v1").Subrouter()
	app.sysRouter = sys_router.NewSystemRouter(app.db)
	app.sysRouter.RegisterRouter(sysSubrouter)

	if err := http.ListenAndServe(fmt.Sprintf(":%v", cfg.HttpServerPort), router); err != nil {
		foree_logger.Logger.Error("Service initial error", "cause", err.Error())
		return err
	}

	return nil
}
