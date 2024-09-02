package logger

type Logger interface {
	Debug(message string, kv ...any)
	Info(message string, kv ...any)
	Warn(message string, kv ...any)
	Error(message string, kv ...any)
	Fatal(message string, kv ...any)
	Debugf(format string, v ...any)
	Infof(format string, v ...any)
	Warnf(format string, v ...any)
	Errorf(format string, v ...any)
	Fatalf(format string, v ...any)
}
