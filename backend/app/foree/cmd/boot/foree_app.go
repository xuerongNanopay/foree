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

	// Initial Repo
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

	// app.transactionService = service.NewTransactionService(
	// 	db,
	// 	app.authService,
	// 	app.userGroupRepo,
	// 	app.foreeTxRepo,
	// 	app.txSummaryRepo,
	// 	app.txQuoteRepo,
	// 	app.rateRepo,
	// 	app.rewardRepo,
	// 	app.dailyTxLimitRepo,
	// 	app.feeRepo,
	// 	app.contactAccountRepo,
	// 	app.interacAccountRepo,
	// 	app.feeJointRepo,
	// )

	//Initial handler

	if err := http.ListenAndServe(fmt.Sprintf(":%v", cfg.HttpServerPort), nil); err != nil {
		return err
	}

	return nil
}
