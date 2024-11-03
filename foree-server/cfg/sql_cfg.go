package cfg

import (
	"context"
	"fmt"
	"strconv"
	"sync"
	"sync/atomic"
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
	cfg, err := c.loadCfg(name, func(conf *configuration) any {
		v := new(string)
		*v = conf.RawValue
		return StringConfig{
			v: v,
		}
	})

	if err != nil {
		return StringConfig{}, err
	}

	return cfg.(StringConfig), nil
}

func (c *SQLCFG) LoadIntCfg(name string) (IntConfig, error) {
	cfg, err := c.loadCfg(name, func(conf *configuration) any {
		v := new(int32)
		i, err := strconv.Atoi(conf.RawValue)
		if err != nil {
			return err
		}
		atomic.StoreInt32(v, int32(i))
		return IntConfig{
			v: v,
		}
	})

	if err != nil {
		return IntConfig{}, err
	}

	return cfg.(IntConfig), nil
}

func (c *SQLCFG) LoadInt64Cfg(name string) (Int64Config, error) {
	cfg, err := c.loadCfg(name, func(conf *configuration) any {
		v := new(int64)
		i, err := strconv.ParseInt(conf.RawValue, 10, 64)
		if err != nil {
			return err
		}
		atomic.StoreInt64(v, i)
		return Int64Config{
			v: v,
		}
	})

	if err != nil {
		return Int64Config{}, err
	}

	return cfg.(Int64Config), nil
}

func (c *SQLCFG) LoadBoolCfg(name string) (BoolConfig, error) {
	cfg, err := c.loadCfg(name, func(conf *configuration) any {
		v := new(uint32)
		i, err := strconv.ParseBool(conf.RawValue)
		if err != nil {
			return err
		}
		if i {
			atomic.StoreUint32(v, 1)
		} else {
			atomic.StoreUint32(v, 0)
		}
		return BoolConfig{
			v: v,
		}
	})

	if err != nil {
		return BoolConfig{}, err
	}

	return cfg.(BoolConfig), nil
}
