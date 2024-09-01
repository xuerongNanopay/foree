package foree_boot

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/go-sql-driver/mysql"
	"xue.io/go-pay/app/foree/account"
	foree_auth "xue.io/go-pay/app/foree/auth"
	foree_config "xue.io/go-pay/app/foree/cmd/config"
	"xue.io/go-pay/app/foree/referral"
	"xue.io/go-pay/app/foree/service"
	"xue.io/go-pay/app/foree/transaction"
	"xue.io/go-pay/auth"
	"xue.io/go-pay/config"
	ms "xue.io/go-pay/db/mysql"
	"xue.io/go-pay/partner/idm"
	"xue.io/go-pay/partner/nbp"
	"xue.io/go-pay/partner/scotia"
)

type ForeeApp struct {
	envFilePath            string
	db                     *sql.DB
	userRepo               *auth.UserRepo
	userGroupRepo          *auth.UserGroupRepo
	sessionRepo            *auth.SessionRepo
	emailPasswdRepo        *auth.EmailPasswdRepo
	rolePermissionRepo     *auth.RolePermissionRepo
	contactAccountRepo     *account.ContactAccountRepo
	interacAccountRepo     *account.InteracAccountRepo
	userExtraRepo          *foree_auth.UserExtraRepo
	userIdnetificationRepo *foree_auth.UserIdentificationRepo
	referralRepo           *referral.ReferralRepo
	dailyTxLimitRepo       *transaction.DailyTxLimitRepo
	feeRepo                *transaction.FeeRepo
	feeJointRepo           *transaction.FeeJointRepo
	rewardRepo             *transaction.RewardRepo
	rateRepo               *transaction.RateRepo
	idmTxRepo              *transaction.IdmTxRepo
	idmRepo                *transaction.IDMComplianceRepo
	foreeTxRepo            *transaction.ForeeTxRepo
	txHistoryRepo          *transaction.TxHistoryRepo
	interacCITxRepo        *transaction.InteracCITxRepo
	interacRefundTxRepo    *transaction.InteracRefundTxRepo
	nbpCOTxRepo            *transaction.NBPCOTxRepo
	txQuoteRepo            *transaction.TxQuoteRepo
	txSummaryRepo          *transaction.TxSummaryRepo
	authService            *service.AuthService
	accountService         *service.AccountService
	transactionService     *service.TransactionService
	scotiaClient           scotia.ScotiaClient
	idmClient              idm.IDMClient
	nbpClient              nbp.NBPClient
	txProcessor            *service.TxProcessor
	ciTxProcessor          *service.CITxProcessor
	idmTxProcessor         *service.IDMTxProcessor
	nbpTxProcessor         *service.NBPTxProcessor
}

func (app *ForeeApp) Boot(envFilePath string) error {
	app.envFilePath = envFilePath
	var cfg foree_config.ForeeLocalConfig
	if err := config.LoadFromFile(&cfg, envFilePath); err != nil {
		return err
	}

	//Initial DB
	db, err := ms.NewMysqlPool(mysql.Config{
		Addr:                 fmt.Sprintf("%s:%s", cfg.MysqlDBHost, cfg.MysqlDBPort),
		DBName:               cfg.MysqlDBName,
		User:                 cfg.MysqlDBUser,
		Passwd:               cfg.MysqlDBPasswd,
		AllowNativePasswords: true,
		ParseTime:            true,
	}, 40, 20)

	if err != nil {
		return err
	}
	app.db = db

	// Initial Repos
	app.userRepo = auth.NewUserRepo(db)
	app.userGroupRepo = auth.NewUserGroupRepo(db)
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
	app.rewardRepo = transaction.NewRewardRepo(db)
	app.idmTxRepo = transaction.NewIdmTxRepo(db)
	app.idmRepo = transaction.NewIDMComplianceRepo(db)
	app.foreeTxRepo = transaction.NewForeeTxRepo(db)
	app.txHistoryRepo = transaction.NewTxHistoryRepo(db)
	app.interacCITxRepo = transaction.NewInteracCITxRepo(db)
	app.interacRefundTxRepo = transaction.NewInteracRefundTxRepo(db)
	app.nbpCOTxRepo = transaction.NewNBPCOTxRepo(db)
	app.txQuoteRepo = transaction.NewTxQuoteRepo(5, 2048)
	app.txSummaryRepo = transaction.NewTxSummaryRepo(db)

	//Initial vendors
	app.scotiaClient = scotia.NewMockScotiaClient()
	app.idmClient = idm.NewMockIDMClient()
	app.nbpClient = nbp.NewMockNBPClient()

	//Initial transaction processors
	app.txProcessor = service.NewTxProcessor(
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
	app.ciTxProcessor = service.NewCITxProcessor(
		app.db,
		service.ScotiaProfile{},
		app.scotiaClient,
		app.interacCITxRepo,
		app.foreeTxRepo,
		app.txSummaryRepo,
		app.txProcessor,
	)
	app.idmTxProcessor = service.NewIDMTxProcessor(
		app.db,
		app.foreeTxRepo,
		app.idmTxRepo,
		app.idmClient,
	)
	app.nbpTxProcessor = service.NewNBPTxProcessor(
		app.db,
		app.foreeTxRepo,
		app.txProcessor,
		app.nbpCOTxRepo,
		app.nbpClient,
		app.userExtraRepo,
		app.userIdnetificationRepo,
	)

	app.txProcessor.SetCITxProcessor(app.ciTxProcessor)
	app.txProcessor.SetIDMTxProcessor(app.idmTxProcessor)
	app.txProcessor.SetNBPTxProcessor(app.nbpTxProcessor)

	if err := app.ciTxProcessor.Start(); err != nil {
		return err
	}

	if err := app.nbpTxProcessor.Start(); err != nil {
		return err
	}

	//Initial service
	app.authService = service.NewAuthService(
		db, app.sessionRepo,
		app.userRepo,
		app.emailPasswdRepo,
		app.rolePermissionRepo,
		app.userIdnetificationRepo,
		app.interacAccountRepo,
		app.userGroupRepo,
	)

	app.accountService = service.NewAccountService(
		app.authService,
		app.contactAccountRepo,
		app.interacAccountRepo,
	)

	app.transactionService = service.NewTransactionService(
		db,
		app.authService,
		app.userGroupRepo,
		app.foreeTxRepo,
		app.txSummaryRepo,
		app.txQuoteRepo,
		app.rateRepo,
		app.rewardRepo,
		app.dailyTxLimitRepo,
		app.feeRepo,
		app.contactAccountRepo,
		app.interacAccountRepo,
		app.feeJointRepo,
		app.txProcessor,
		app.scotiaClient,
		app.nbpClient,
	)

	//Initial handler

	if err := http.ListenAndServe(fmt.Sprintf(":%v", cfg.HttpServerPort), nil); err != nil {
		return err
	}

	return nil
}
