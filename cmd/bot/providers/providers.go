package providers

import (
	"os"
	"weatherTGBot/cmd/bot/config"
	"weatherTGBot/pkg/db"
	"weatherTGBot/pkg/db/postgres"
	"weatherTGBot/pkg/logger"
	"weatherTGBot/pkg/logger/zerolog"
	"weatherTGBot/pkg/telegram"
)

func ProvideConsoleLogger(cnf config.Config) logger.Logger {
	return zerolog.NewZeroLog(os.Stderr, cnf.Logger.Lvl)
}

func ProvideFileLogger(cnf config.Config) (logger.Logger, func(), error) {
	var ioWriter = os.Stdout
	var closeFn = func() {}

	if cnf.Logger.FilePath != "" {
		logfile, err := os.OpenFile(cnf.Logger.FilePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0775)
		if err != nil {
			return nil, nil, err
		}

		closeFn = func() {
			_ = logfile.Close()
		}

		ioWriter = logfile
	}

	return zerolog.NewZeroLog(ioWriter, cnf.Logger.Lvl), closeFn, nil
}

func ProvideTgBotRepo(cnf config.Config, log logger.Logger) (db.TgBotRepo, func(), error) {
	var closeFn = func() {}

	repo, err := postgres.NewDBConnection(config.GetPgDsn(cnf), log)
	if err != nil {
		return nil, nil, err
	}

	closeFn = func() {
		_ = repo.Close()
	}

	return repo, closeFn, nil
}

func ProvideWeatherApi(cnf config.Config) (telegram.WeatherApi, error) {
	weatherApi, err := telegram.NewOpenWeatherMapApi(cnf.WeatherApi.Unit, cnf.WeatherApi.Lang, cnf.WeatherApi.APIKey)
	if err != nil {
		return nil, err
	}

	return weatherApi, nil
}
