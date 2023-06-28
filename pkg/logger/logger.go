package logger

type Logger interface {
	Warn(kv ...interface{})
	Error(kv ...interface{})
	Debug(kv ...interface{})
	Info(kv ...interface{})

	Warnf(s string, kv ...interface{})
	Errorf(s string, kv ...interface{})
	Debugf(s string, kv ...interface{})
	Infof(s string, kv ...interface{})
}
