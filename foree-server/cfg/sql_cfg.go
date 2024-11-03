package cfg

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type SQLCFG struct {
	mu      sync.Mutex
	configs sync.Map
	repo    *configurationRepo
}

type configWrapper struct {
	config    any
	rawValue  string
	expiredAt time.Time
}

func (c *SQLCFG) loadCfg(name string, converter func(conf *configuration) any) (any, error) {
	if cfg, ok := c.configs.Load(name); ok {
		w := cfg.(configWrapper)
		return w.config, nil
	}
	c.mu.Lock()
	defer c.mu.Unlock()

	if cfg, ok := c.configs.Load(name); ok {
		w := cfg.(configWrapper)
		return w.config, nil
	}

	conf, err := c.repo.getUniqueConfigurationByName(context.TODO(), name)
	if err != nil {
		return name, err
	}
	if conf == nil {
		return name, fmt.Errorf("configuraion `%v` not found", name)
	}

	cw := configWrapper{
		rawValue:  conf.RawValue,
		expiredAt: time.Now().Add(time.Millisecond * time.Duration(conf.RefreshInterval)),
		config:    converter(conf),
	}

	c.configs.Store(name, cw)
	return cw.config, nil
}

func (c *SQLCFG) LoadStringCfg(name string) (StringConfig, error) {
	if cfg, ok := c.configs.Load(name); ok {
		w := cfg.(configWrapper)
		return w.config.(StringConfig), nil
	}
	c.mu.Lock()
	defer c.mu.Unlock()

	if cfg, ok := c.configs.Load(name); ok {
		w := cfg.(configWrapper)
		return w.config.(StringConfig), nil
	}

	conf, err := c.repo.getUniqueConfigurationByName(context.TODO(), name)
	if err != nil {
		return StringConfig{}, err
	}
	if conf == nil {
		return StringConfig{}, fmt.Errorf("configuraion `%v` not found", name)
	}

	v := conf.RawValue
	cw := configWrapper{
		rawValue:  conf.RawValue,
		expiredAt: time.Now().Add(time.Millisecond * time.Duration(conf.RefreshInterval)),
		config: StringConfig{
			v: &v,
		},
	}

	c.configs.Store(name, cw)
	return cw.config.(StringConfig), nil
}
