package scotia

import (
	"testing"
)

func TestScotiaConfig(t *testing.T) {

	config := map[string]string{
		ConfigBaseUrl:           "http://www.dummy.com",
		ConfigBasicAuthUsername: "xue",
		ConfigBasicAuthPassword: "11111",
		ConfigClientId:          "xxuuee",
		ConfigJWTKid:            "yyyy",
		ConfigJWTAudience:       "zzzz",
		ConfigJWTExpiryMinutes:  "300",
		ConfigPrivateKeyDir:     "~/test",
		ConfigPublicKeyDir:      "~/test_pub",
		ConfigScope:             "lllll",
	}

	t.Run("config should construct correctly", func(t *testing.T) {

		sc := NewScotiaConfigWithDefaultConfig(config)

		show := sc.ShowConfigs()

		for key, value := range config {
			if value != show[key] {
				t.Errorf("expect %v, but %v", config[key], value)
			}
		}

	})

	t.Run("config should be a reference type", func(t *testing.T) {

		sc := NewScotiaConfigWithDefaultConfig(config)

		sc1 := sc
		sc.SetScope("new scope")
		if sc1.GetScope() != sc.GetScope() {
			t.Errorf("expect scotiaConfig shoule be a reference type")
		}

	})

	t.Run("SetConfig should work", func(t *testing.T) {

		sc := NewScotiaConfigWithDefaultConfig(config)

		sc.SetConfig(ConfigBaseUrl, "aaa")

		if sc.GetBaseUrl() != "aaa" {
			t.Errorf("expect %v, but %v", "aaa", sc.GetBaseUrl())
		}

	})
}
