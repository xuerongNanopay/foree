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
		rawTag, ok := sTag.Lookup("env_var")
		if ok {
			tags := strings.Split(rawTag, ",")
			if len(tags) != 2 {
				return fmt.Errorf("invalid env_var format `%s` in field `%s`", rawTag, sField.Name)
			}
			label := tags[0]
			l2 := tags[1]
			if l2 != "required" && !strings.HasPrefix(l2, "default=") {
				return fmt.Errorf("invalid env_var format `%s` in field `%s`", rawTag, sField.Name)
			}

			if l2 == "required" {
				value, ok := os.LookupEnv(label)
				if !ok {
					return fmt.Errorf("do not find `%s` in environment", label)
				}
				err := reflect_util.SetStuctValueFromString(config, sField.Name, value)
				if err != nil {
					return err
				}
			} else {
				value, ok := os.LookupEnv(label)
				if !ok || value == "" {
					if len(strings.Split(l2, "=")) != 2 {
						return fmt.Errorf("invalid env_var format `%s` in field `%s`", rawTag, sField.Name)
					}
					value = strings.Split(l2, "=")[1]
				}
				err := reflect_util.SetStuctValueFromString(config, sField.Name, value)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}
