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
			label := tags[0]
			fmt.Println("label: ", label)
			l2 := tags[1]
			if l2 == "required" {
				value, ok := os.LookupEnv(sField.Name)
				if !ok {
					return fmt.Errorf("do not find `%s` in environment", sField.Name)
				}
				reflect_util.SetStuctValueFromString(config, sField.Name, value)
			} else {
				value, ok := os.LookupEnv(sField.Name)
				if !ok {
					return fmt.Errorf("do not find `%s` in environment", sField.Name)
				}
				if value == "" {
					value = strings.Split(l2, "=")[1]
				}
				reflect_util.SetStuctValueFromString(config, sField.Name, value)
			}
		}
	}
	return nil
}
