package cfg

type CFG interface {
	LoadCfg(key string) (any, bool)
	LoadStringCfg(key string) (string, bool)
	LoadIntCfg(key string) (int, bool)
	LoadInt64Cfg(key string) (int64, bool)
}
