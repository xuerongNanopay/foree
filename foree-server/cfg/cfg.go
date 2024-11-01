package cfg

import "database/sql"

type CFG interface {
	LoadStringCfg(name string) (StringConfig, error)
	LoadBoolCfgBool(name string) (BoolConfig, error)
	LoadIntCfg(name string) (IntConfig, error)
	LoadInt64Cfg(name string) (Int64Config, error)
	Reset(name string)
	ResetAll()
}

type SQLConfigure struct {
	db      *sql.DB
	configs map[string]Config[any]
}
