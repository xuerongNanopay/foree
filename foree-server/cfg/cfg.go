package cfg

type CFG interface {
	LoadStringCfg(name string) (StringConfig, error)
	LoadBoolCfg(name string) (BoolConfig, error)
	LoadIntCfg(name string) (IntConfig, error)
	LoadInt64Cfg(name string) (Int64Config, error)
	Reset(name string)
	ResetAll()
}
