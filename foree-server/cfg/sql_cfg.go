package cfg

type SQLCFG struct {
	configurationRepo *configurationRepo
	configs           map[string]Config[any]
}
