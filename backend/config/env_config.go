package config

import (
	"fmt"
	"os"
	"strings"

	reflect_util "xue.io/go-pay/util/reflect"
)

func Load(config any) error {
	for _, f := range reflect_util.GetAllFieldNamesOfStruct(config) {
		sField, sTag := reflect_util.GetTagOfStruct(config, f)
		rawTag, ok := sTag.Lookup("os_config")
		if ok {
			tags := strings.Split(rawTag, ",")
			if len(tags) != 2 {
				return fmt.Errorf("invalid os_config format `%s` in field `%s`", tags, sField.Name)
			}
			label := tags[0]
			l2 := tags[1]
			if l2 != "required" || !strings.HasPrefix(l2, "default=") {
				return fmt.Errorf("invalid os_config format `%s` in field `%s`", tags, sField.Name)
			}

			if l2 == "required" {
				value, ok := os.LookupEnv(label)
				if !ok {
					return fmt.Errorf("do not find `%s` in environment", label)
				}
				return reflect_util.SetStuctValueFromString(config, sField.Name, value)
			} else {
				value, ok := os.LookupEnv(label)
				if !ok || value == "" {
					value = strings.Split(l2, "=")[1]
				}
				return reflect_util.SetStuctValueFromString(config, sField.Name, value)
			}
		}
	}
	return nil
}
