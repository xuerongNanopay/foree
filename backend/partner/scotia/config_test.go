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
		ConfigJWTExpiry:         "300",
		ConfigPrivateKeyDir:     "~/test",
		ConfigPublicKeyDir:      "~/test_pub",
		ConfigScope:             "lllll",
	}

	t.Run("config should construct correctly", func(t *testing.T) {

		sc := NewScotiaConfigWithDefaultConfig(config)

		if config[ConfigBaseUrl] != sc.GetBaseUrl() {
			t.Errorf("expect %v, but %v", config[ConfigBaseUrl], sc.GetBaseUrl())
		}

		if sc.GetJWTExpiry() != 300 {
			t.Errorf("expect %v, but %v", 600, sc.GetJWTExpiry())
		}

		show := sc.ShowConfigs()

		for key, value := range show {
			if value != config[key] {
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
}
