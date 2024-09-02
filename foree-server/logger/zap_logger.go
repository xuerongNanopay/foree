package logger

import (
	"encoding/json"
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewZapLogger(level, outputPath string) (*ZapLogger, error) {
	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.TimeKey = "timestamp"
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderCfg.MessageKey = "message"
	encoderCfg.LevelKey = "level"
	encoderCfg.EncodeLevel = zapcore.CapitalLevelEncoder
	//TODO: custom empty encode caller
	encoderCfg.EncodeCaller = zapcore.ShortCallerEncoder

	rawJSON := []byte(fmt.Sprintf(`{
		"level": "%s",
		"encoding": "json",
		"outputPaths": ["stdout", "%s"],
		"errorOutputPaths": ["stderr", "%s"]
	}`, level, outputPath, outputPath))

	var cfg zap.Config
	if err := json.Unmarshal(rawJSON, &cfg); err != nil {
		return nil, err
	}

	cfg.EncoderConfig = encoderCfg

	logger := zap.Must(cfg.Build()).Sugar()

	return &ZapLogger{
		log: logger,
	}, nil
}

type ZapLogger struct {
	log *zap.SugaredLogger
}

func (z *ZapLogger) Debug(message string, kv ...any) {
	defer z.log.Sync()
	z.log.Debugw(message, kv...)
}
func (z *ZapLogger) Info(message string, kv ...any) {
	defer z.log.Sync()
	z.log.Infow(message, kv...)
}
func (z *ZapLogger) Warn(message string, kv ...any) {
	defer z.log.Sync()
	z.log.Warnw(message, kv...)
}
func (z *ZapLogger) Error(message string, kv ...any) {
	defer z.log.Sync()
	z.log.Errorw(message, kv...)
}
func (z *ZapLogger) Fatal(message string, kv ...any) {
	defer z.log.Sync()
	z.log.Fatalw(message, kv)
}

func (z *ZapLogger) Debugf(format string, v ...any) {
	defer z.log.Sync()
	z.log.Debugf(format, v...)
}
func (z *ZapLogger) Infof(format string, v ...any) {
	defer z.log.Sync()
	z.log.Infof(format, v...)
}
func (z *ZapLogger) Warnf(format string, v ...any) {
	defer z.log.Sync()
	z.log.Warnf(format, v...)
}
func (z *ZapLogger) Errorf(format string, v ...any) {
	defer z.log.Sync()
	z.log.Errorf(format, v...)
}
func (z *ZapLogger) Fatalf(format string, v ...any) {
	defer z.log.Sync()
	z.log.Fatalf(format, v...)
}
