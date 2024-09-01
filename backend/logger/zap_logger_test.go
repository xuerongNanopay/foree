package logger

import (
	"testing"
	"time"

	"go.uber.org/zap"
)

func TestZap(t *testing.T) {

	t.Run("demo 1", func(t *testing.T) {
		sugar := zap.NewExample().Sugar()
		defer sugar.Sync()
		sugar.Infow("failed to fetch URL",
			"url", "http://example.com",
			"attempt", 3,
			"backoff", time.Second,
		)
		sugar.Infof("failed to fetch URL: %s", "http://example.com")
	})

	t.Run("ZapLogger demo", func(t *testing.T) {
		logger, err := NewZapLogger("debug", "/tmp/zap_out")
		if err != nil {
			t.Fatal(err.Error())
		}
		logger.Debug("Mesage!!!", "myKey", "mayValue")
		logger.Errorf("EEEEEEEEEEE")
	})
}
