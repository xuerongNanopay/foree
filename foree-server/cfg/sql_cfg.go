package cfg

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"xue.io/go-pay/logger"
)

const refreshInterval = 5 * time.Minute

func NewSQLCFG(db *sql.DB, logger logger.Logger) *SQLCFG {
	ret := &SQLCFG{
		refreshTicker: time.NewTicker(refreshInterval),
		logger:        logger,
		repo: &configurationRepo{
			db: db,
		},
	}

	go ret.startRefresher()
	return ret
}

type SQLCFG struct {
	mu            sync.Mutex
	configs       sync.Map
	hardRefresh   chan struct{}
	refreshTicker *time.Ticker
	repo          *configurationRepo
	logger        logger.Logger
}

type configWrapper struct {
	config    any
	rawValue  string
	expiredAt time.Time
}

func (c *SQLCFG) startRefresher() {
	for {
		select {
		case <-c.hardRefresh:
			c.refresh()
		case <-c.refreshTicker.C:
			c.refresh()
		}
	}
}

func (c *SQLCFG) refresh() {
	c.mu.Lock()
	defer c.mu.Unlock()

	names := make([]string, 0)
	c.configs.Range(func(k, v interface{}) bool {
		name := k.(string)
		cw := v.(configWrapper)
		if cw.expiredAt.Before(time.Now()) {
			names = append(names, name)
		}
		return true
	})
	confs, err := c.repo.getAllConfigurationByNames(context.TODO(), names...)
	if err != nil {
		c.logger.Error("SQLCFG_Refresh_FAIL", "cause", err)
	}

	for _, newConf := range confs {
		v, ok := c.configs.Load(newConf.Name)
		if !ok {
			c.logger.Error("SQLCFG_Refresh_FAIL", "name", newConf.Name, "cause", "configuration not found")
		}
		cw := v.(configWrapper)

		if cw.rawValue == newConf.RawValue {
			continue
		}

		nCw := cw
		nCw.rawValue = newConf.RawValue

		expiry := time.Now()
		if newConf.RefreshInterval < 0 {
			expiry = expiry.Add(time.Hour * time.Duration(24*356))
		} else {
			expiry = expiry.Add(time.Millisecond * time.Duration(newConf.RefreshInterval))
		}

		nCw.expiredAt = expiry

		switch curConf := nCw.config.(type) {
		case StringConfig:
			curConf.v.Swap(newConf.RawValue)
		case IntConfig:
			i, err := strconv.Atoi(newConf.RawValue)
			if err != nil {
				c.logger.Error("SQLCFG_Refresh_FAIL", "cause", err)
			} else {
				atomic.StoreInt32(curConf.v, int32(i))
			}
		case Int64Config:
			i, err := strconv.ParseInt(newConf.RawValue, 10, 64)
			if err != nil {
				c.logger.Error("SQLCFG_Refresh_FAIL", "cause", err)
			} else {
				atomic.StoreInt64(curConf.v, i)
			}
		case BoolConfig:
			i, err := strconv.ParseBool(newConf.RawValue)
			if err != nil {
				c.logger.Error("SQLCFG_Refresh_FAIL", "cause", err)
			} else {
				if i {
					atomic.StoreUint32(curConf.v, 1)
				} else {
					atomic.StoreUint32(curConf.v, 0)
				}
			}
		default:
			c.logger.Error("SQLCFG_Refresh_FAIL", "dataType", fmt.Sprintf("%T", curConf), "cause", "unknown config type")
			continue
		}
		c.configs.Swap(newConf.Name, nCw)
		c.logger.Info("SQLCFG_Refresh_SUCCESS", "configurationName", newConf.Name)
	}
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

	expiry := time.Now()
	if conf.RefreshInterval < 0 {
		expiry = expiry.Add(time.Hour * time.Duration(24*356))
	} else {
		expiry = expiry.Add(time.Millisecond * time.Duration(conf.RefreshInterval))
	}

	cw := configWrapper{
		rawValue:  conf.RawValue,
		expiredAt: expiry,
		config:    converter(conf),
	}

	c.configs.Store(name, cw)
	return cw.config, nil
}

func (c *SQLCFG) LoadStringCfg(name string) (StringConfig, error) {
	cfg, err := c.loadCfg(name, func(conf *configuration) any {
		v := new(atomic.Value)
		v.Store(conf.RawValue)
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

func (c *SQLCFG) Reset(name string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	v, ok := c.configs.Load(name)
	if !ok {
		return
	}

	cw := v.(configWrapper)

	var zero time.Time
	nCw := cw
	nCw.expiredAt = zero

	c.configs.Swap(name, nCw)
	select {
	case c.hardRefresh <- struct{}{}:
	default:
	}
}

func (c *SQLCFG) ResetAll() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.configs.Range(func(k, v interface{}) bool {
		name := k.(string)
		cw := v.(configWrapper)

		var zero time.Time
		nCw := cw
		nCw.expiredAt = zero

		c.configs.Swap(name, nCw)
		return true
	})

	select {
	case c.hardRefresh <- struct{}{}:
	default:
	}
}
