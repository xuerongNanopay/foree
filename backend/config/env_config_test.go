package config

import (
	"fmt"
	"testing"
)

func TestEnvConfig(t *testing.T) {
	type Person struct {
		Name string `os_config:"name,required,default=xue"`
		Age  int    `os_config:"age,required"`
	}

	var person Person

	fmt.Println(person)
}
