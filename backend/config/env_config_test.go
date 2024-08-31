package config

import (
	"os"
	"testing"
)

func TestEnvConfig(t *testing.T) {

	t.Run("environment config should load", func(t *testing.T) {
		type Person struct {
			FirstName  string `os_config:"FIRST_NAME,required"`
			MiddleName string `os_config:"MIDDLE_NAME,default=rong"`
			Age        int    `os_config:"AGE,required"`
			Male       bool   `os_config:"MALE,required"`
		}

		var person Person

		os.Setenv("FIRST_NAME", "xue")
		os.Setenv("AGE", "35")
		os.Setenv("MALE", "true")

		err := Load(&person)

		if err != nil {
			t.Errorf("should not raise error `%v`", err.Error())
		}

		if person.FirstName != "xue" &&
			person.MiddleName != "rong" &&
			person.Age != 35 &&
			person.Male != true {

			t.Errorf("should load environment successfully, but got `%v`", person)
		}
	})
}
