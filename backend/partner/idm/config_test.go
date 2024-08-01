package idm

import (
	"testing"
)

func TestIDMConfig(t *testing.T) {

	config := map[string]string{
		ConfigBaseUrl:           "http://www.dummy.com",
		ConfigBasicAuthUsername: "xue",
		ConfigBasicAuthPassword: "11111",
		ConfigHashingSalt:       "xxuuee",
		ConfigProfile:           "yyyy",
	}

	t.Run("config should construct correctly", func(t *testing.T) {

		sc := NewIDMConfigWithDefaultConfig(config)

		show := sc.ShowConfigs()

		for key, value := range show {
			if value != config[key] {
				t.Errorf("expect %v, but %v", config[key], value)
			}
		}
	})

	t.Run("config should be a reference type", func(t *testing.T) {

		sc := NewIDMConfigWithDefaultConfig(config)

		sc1 := sc
		sc.SetProfile("new profile")
		if sc1.GetProfile() != sc.GetProfile() {
			t.Errorf("expect scotiaConfig shoule be a reference type")
		}

	})

	t.Run("SetConfig should work", func(t *testing.T) {

		sc := NewIDMConfigWithDefaultConfig(config)

		sc.SetConfig(ConfigBaseUrl, "aaa")

		if sc.GetBaseUrl() != "aaa" {
			t.Errorf("expect %v, but %v", "aaa", sc.GetBaseUrl())
		}

	})
}
