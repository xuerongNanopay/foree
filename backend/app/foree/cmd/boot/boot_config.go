package foree_boot

import (
	"strconv"
)

type BootConfig struct {
	MysqlDBHost    string
	MysqlDBAddress string
	MysqlDBUser    string
	MysqlDBPasswd  string
	HttpServerHost string
	HttpServerPort string
}

// func (b *BootConfig) load(envPath string) error {

// }

// func NewBootConfig(envPath string) (*BootConfig, error) {
// 	ex, err := os.Executable()
// 	if err != nil {
// 		return nil, err
// 	}
// 	return nil
// }

type BootConfigField[V comparable] struct {
	Key     string
	Require bool
	Default V
	Parser  func(r string) V
}

func (b BootConfigField[V]) IsRequire() bool {
	return b.Require
}

func (b BootConfigField[V]) GetValue(r string) V {
	v := b.Parser(r)
	if b.isZero(v) {
		return b.Default
	}
	return v
}

func (b BootConfigField[V]) isZero(v V) bool {
	var z V
	return v == z
}

func NewIntConfigFieldWithDefault(key string, d int) BootConfigField[int] {
	return BootConfigField[int]{
		Key:     key,
		Require: false,
		Default: d,
		Parser: func(r string) int {
			i, err := strconv.Atoi(r)
			if err != nil {
				return 0
			}
			return i
		},
	}
}

func NewIntConfigField(key string) BootConfigField[int] {
	return BootConfigField[int]{
		Key:     key,
		Require: true,
		Parser: func(r string) int {
			i, err := strconv.Atoi(r)
			if err != nil {
				return 0
			}
			return i
		},
	}
}

func NewStringConfigFieldWithDefault(key string, d string) BootConfigField[string] {
	return BootConfigField[string]{
		Key:     key,
		Require: false,
		Default: d,
		Parser: func(r string) string {
			return r
		},
	}
}

func NewStringConfigField(key string) BootConfigField[string] {
	return BootConfigField[string]{
		Key:     key,
		Require: true,
		Parser: func(r string) string {
			return r
		},
	}
}

func NewBoolConfigFieldWithDefault(key string, d bool) BootConfigField[bool] {
	return BootConfigField[bool]{
		Key:     key,
		Require: false,
		Default: d,
		Parser: func(r string) bool {
			i, err := strconv.ParseBool(r)
			if err != nil {
				return false
			}
			return i
		},
	}
}

func NewBoolConfigField(key string) BootConfigField[bool] {
	return BootConfigField[bool]{
		Key:     key,
		Require: true,
		Parser: func(r string) bool {
			i, err := strconv.ParseBool(r)
			if err != nil {
				return false
			}
			return i
		},
	}
}
