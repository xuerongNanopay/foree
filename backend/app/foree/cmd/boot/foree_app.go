package foree_boot

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/go-sql-driver/mysql"
	foree_config "xue.io/go-pay/app/foree/cmd/config"
	"xue.io/go-pay/auth"
	"xue.io/go-pay/config"
	ms "xue.io/go-pay/db/mysql"
)

type ForeeApp struct {
	envFilePath     string
	db              *sql.DB
	userRepo        *auth.UserRepo
	userGroupRepo   *auth.UserGroupRepo
	sessionRepo     *auth.SessionRepo
	emailPasswdRepo *auth.EmailPasswdRepo
}

func (app *ForeeApp) boot(envFilePath string) error {
	app.envFilePath = envFilePath
	var cfg foree_config.ForeeLocalConfig
	if err := config.LoadFromFile(&cfg, envFilePath); err != nil {
		return err
	}

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

	app.userRepo = auth.NewUserRepo(db)
	app.userGroupRepo = auth.NewUserGroupRepo(db)
	app.sessionRepo = auth.NewDefaultSessionRepo(db)
	app.emailPasswdRepo = auth.NewEmailPasswdRepo(db)

	//Initial DB
	//Initial Repo
	//Initial service
	//Initial handler

	if err := http.ListenAndServe(fmt.Sprintf(":%v", cfg.HttpServerPort), nil); err != nil {
		return err
	}

	return nil
}
