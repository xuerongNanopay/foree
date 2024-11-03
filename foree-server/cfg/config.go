package cfg

type Config[T any] interface {
	Get() T
}

type IntConfig struct {
	v *int
}

func (c *IntConfig) Get() int {
	return *c.v
}

type BoolConfig struct {
	v *bool
}

func (c *BoolConfig) Get() bool {
	return *c.v
}

type Int64Config struct {
	v *int64
}

func (c *Int64Config) Get() int64 {
	return *c.v
}

type StringConfig struct {
	v *string
}

func (c *StringConfig) Get() string {
	return *c.v
}
