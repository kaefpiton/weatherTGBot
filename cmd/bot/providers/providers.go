package providers

import (
	"os"
	"weatherTGBot/cmd/bot/config"
	"weatherTGBot/pkg/db"
	"weatherTGBot/pkg/db/postgres"
	"weatherTGBot/pkg/logger"
)

// todo сделать через конфиг
func ProvideConsoleLogger(logLvl string) logger.Logger {
	return logger.NewZeroLog(os.Stderr, logLvl)
}

// todo сделать через конфиг
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

func ProvideTgBotRepo(cnf config.Config) (db.TgBotRepo, func(), error) {
	var closeFn = func() {}

	//todo добавить в dsn еще конфигов
	repo, err := postgres.NewDBConnection(config.GetPgDsn(cnf))
	if err != nil {
		return nil, nil, err
	}

	closeFn = func() {
		_ = repo.Close()
	}

	return repo, closeFn, nil
}
