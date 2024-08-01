package nbp

import (
	"testing"
)

func TestNBPConfig(t *testing.T) {

	config := map[string]string{
		ConfigBaseUrl:              "http://www.dummy.com",
		ConfigBasicAuthUsername:    "xue",
		ConfigBasicAuthPassword:    "11111",
		ConfigAgencyCode:           "xxuuee",
		ConfigAuthAttempts:         "111",
		ConfigTokenExpiryThreshold: "222",
	}

	t.Run("config should construct correctly", func(t *testing.T) {

		sc := NewNBPConfigWithDefaultConfig(config)

		show := sc.ShowConfigs()

		if sc.GetAuthAttempts() != 111 {
			t.Errorf("expect %v, but %v", 111, sc.GetAuthAttempts())
		}

		if sc.GetTokenExpiryThreshold() != 222 {
			t.Errorf("expect %v, but %v", 222, sc.GetTokenExpiryThreshold())
		}

		for key, value := range show {
			if value != config[key] {
				t.Errorf("expect %v, but %v", config[key], value)
			}
		}
	})

	t.Run("config should be a reference type", func(t *testing.T) {

		sc := NewNBPConfigWithDefaultConfig(config)

		sc1 := sc
		sc.SetAgencyCode("new Agency Code")
		if sc1.GetAgencyCode() != sc.GetAgencyCode() {
			t.Errorf("expect scotiaConfig shoule be a reference type")
		}

	})

	t.Run("SetConfig should work", func(t *testing.T) {

		sc := NewNBPConfigWithDefaultConfig(config)

		sc.SetConfig(ConfigBaseUrl, "aaa")

		if sc.GetBaseUrl() != "aaa" {
			t.Errorf("expect %v, but %v", "aaa", sc.GetBaseUrl())
		}

	})
}
