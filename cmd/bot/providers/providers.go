package providers

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"os"
	"weatherTGBot/internal/config"
	"weatherTGBot/internal/infrastructure/repository"
	"weatherTGBot/pkg/db/postgres"
	"weatherTGBot/pkg/logger"
	"weatherTGBot/pkg/logger/zerolog"
	"weatherTGBot/pkg/weather"
)

func ProvideConsoleLogger(cnf *config.Config) logger.Logger {
	return zerolog.NewZeroLog(os.Stderr, cnf.Logger.Lvl)
}

func ProvideFileLogger(cnf *config.Config) (logger.Logger, func(), error) {
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

func ProvideDB(cnf *config.Config) (*postgres.DB, func(), error) {
	var closeFn = func() {}

	db, err := postgres.NewDBConnection(cnf)
	if err != nil {
		return nil, nil, err
	}

	closeFn = func() {
		_ = db.Close()
	}

	return db, closeFn, nil
}

func ProvideTgBotRepo(db *postgres.DB, logger logger.Logger) *repository.TgBotRepository {
	return repository.NewBotRepository(db, logger)
}

func ProvideWeatherApi(cnf *config.Config) (weather.WeatherApi, error) {
	weatherApi, err := weather.NewOpenWeatherMapApi(cnf.WeatherApi.Unit, cnf.WeatherApi.Lang, cnf.WeatherApi.APIKey)
	if err != nil {
		return nil, err
	}

	return weatherApi, nil
}

func ProvideBotApi(cnf *config.Config) (*tgbotapi.BotAPI, error) {
	botApi, err := tgbotapi.NewBotAPI(cnf.TelegramApi.APIKey)
	if err != nil {
		return nil, err
	}
	botApi.Debug = cnf.TelegramApi.Debug

	return botApi, err
}
