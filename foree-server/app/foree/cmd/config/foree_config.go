package foree_config

type ForeeLocalConfig struct {
	MysqlDBHost    string `env_var:"MYSQL_DB_HOST,default=localhost"`
	MysqlDBPort    string `env_var:"MYSQL_DB_PORT,default=3306"`
	MysqlDBUser    string `env_var:"MYSQL_DB_USER,default=root"`
	MysqlDBPasswd  string `env_var:"MYSQL_DB_PASSWD,required"`
	MysqlDBName    string `env_var:"MYSQL_DB_NAME,default=foree"`
	HttpServerPort string `env_var:"HTTP_SERVER_PORT,default=8080"`
}
