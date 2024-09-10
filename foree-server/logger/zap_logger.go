package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

func NewZapLogger(level, outputPath string) (*ZapLogger, error) {
	stdout := zapcore.AddSync(os.Stdout)
	file := zapcore.AddSync(&lumberjack.Logger{
		Filename:   outputPath,
		MaxSize:    10, // megabytes
		MaxBackups: 0,
		MaxAge:     0, // days
	})

	l := zap.NewAtomicLevelAt(zap.InfoLevel)

	productionCfg := zap.NewProductionEncoderConfig()
	productionCfg.TimeKey = "timestamp"
	productionCfg.EncodeTime = zapcore.ISO8601TimeEncoder
	productionCfg.MessageKey = "event"
	productionCfg.LevelKey = "level"
	productionCfg.EncodeLevel = zapcore.CapitalLevelEncoder
	//TODO: custom empty encode caller
	productionCfg.EncodeCaller = zapcore.ShortCallerEncoder
	productionCfg.StacktraceKey = "stacktrace"

	developmentCfg := zap.NewDevelopmentEncoderConfig()
	developmentCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder
	developmentCfg.StacktraceKey = "stacktrace"

	consoleEncoder := zapcore.NewConsoleEncoder(developmentCfg)
	fileEncoder := zapcore.NewJSONEncoder(productionCfg)

	core := zapcore.NewTee(
		zapcore.NewCore(consoleEncoder, stdout, l),
		zapcore.NewCore(fileEncoder, file, l),
	)

	return &ZapLogger{
		log: zap.New(core, zap.WithCaller(false), zap.AddStacktrace(zap.ErrorLevel)).Sugar(),
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
