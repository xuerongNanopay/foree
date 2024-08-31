package config

import (
	"fmt"
	"strings"
	"testing"

	reflect_util "xue.io/go-pay/util/reflect"
)

func TestEnvConfig(t *testing.T) {
	type Person struct {
		Name string `os_config:"name,default=xue"`
		Age  int    `os_config:"age,required"`
	}

	var person Person

	for _, f := range reflect_util.GetAllFieldNamesOfStruct(&person) {
		field, t := reflect_util.GetTagOfStruct(&person, f)
		rawTag, ok := t.Lookup("os_config")
		if ok {
			tags := strings.Split(rawTag, ",")
			label := tags[0]
			fmt.Println("label: ", label)
			l2 := tags[1]
			if l2 == "required" {
				//TODO: load from env.
			} else {
				defaultValue := strings.Split(l2, "=")[1]
				reflect_util.TrySetStuctValueFromString(&person, field.Name, defaultValue)
			}
		}
	}
}
