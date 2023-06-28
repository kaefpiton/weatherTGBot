package providers

import (
	"os"
	"weatherTGBot/pkg/logger"
)

func ProvideConsoleLogger(logLvl string) logger.Logger {
	return logger.NewZeroLog(os.Stderr, logLvl)
}

func ProvideFileLogger(logLvl string, filePath string) (logger.Logger, func(), error) {
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

	return logger.NewZeroLog(ioWriter, logLvl), closeFn, nil
}

//todo - добавить dbprovider и остальные провайдеры
