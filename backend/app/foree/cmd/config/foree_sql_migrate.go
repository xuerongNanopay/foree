package foree_config

type ForeeMigrateConfig struct {
	MysqlDBHost   string `env_var:"MYSQL_DB_HOST,default=localhost"`
	MysqlDBPort   string `env_var:"MYSQL_DB_PORT,default=3306"`
	MysqlDBUser   string `env_var:"MYSQL_DB_USER,default=root"`
	MysqlDBPasswd string `env_var:"MYSQL_DB_PASSWD,required"`
	MysqlDBName   string `env_var:"MYSQL_DB_NAME,default=foree"`
}
