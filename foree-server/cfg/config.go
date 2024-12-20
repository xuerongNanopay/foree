package cfg

import "sync/atomic"

type BoolConfig struct {
	v *uint32
}

func (c *BoolConfig) Get() bool {
	return atomic.LoadUint32(c.v) != 0
}

type IntConfig struct {
	v *int32
}

func (c *IntConfig) Get() int {
	return int(atomic.LoadInt32(c.v))
}

type Int64Config struct {
	v *int64
}

func (c *Int64Config) Get() int64 {
	return *c.v
}

type StringConfig struct {
	v *atomic.Value
}

func (c *StringConfig) Get() string {
	return c.v.Load().(string)
}
