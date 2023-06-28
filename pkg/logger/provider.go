package logger

import (
	"os"
)

// todo вынести
func ProvideConsoleLogger(logLvl string) Logger {
	return NewZeroLog(os.Stderr, logLvl)
}

func ProvideFileLogger(logLvl string, filePath string) (Logger, func(), error) {
	var ioWriter = os.Stdout
	var closeFn = func() {}

	if filePath != "" {
		logfile, err := os.OpenFile(filePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0775)
		if err != nil {
			return nil, nil, err
		}

		closeFn = func() {
			_ = logfile.Close()
		}

		ioWriter = logfile
	}

	return NewZeroLog(ioWriter, logLvl), closeFn, nil
}
