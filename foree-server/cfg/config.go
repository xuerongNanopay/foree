package cfg

type loader[T any] func() T
type resetter[T any] func()

type Config[T any] interface {
	Get() T
	Reset()
}

type IntConfig struct {
	l loader[int]
	r resetter[int]
}

func (c *IntConfig) Get() int {
	return c.l()
}

func (c *IntConfig) Reset() {
	c.r()
}

type BoolConfig struct {
	l loader[bool]
	r resetter[bool]
}

func (c *BoolConfig) Get() bool {
	return c.l()
}

func (c *BoolConfig) Reset() {
	c.r()
}

type Int64Config struct {
	l loader[int64]
	r resetter[int64]
}

func (c *Int64Config) Get() int64 {
	return c.l()
}

func (c *Int64Config) Reset() {
	c.r()
}

type StringConfig struct {
	l loader[string]
	r resetter[string]
}

func (c *StringConfig) Get() string {
	return c.l()
}

func (c *StringConfig) Reset() {
	c.r()
}
