package main

import (
	"context"
	"fmt"
	owm "github.com/briandowns/openweathermap"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"os"
	"os/signal"
	"syscall"
	"weatherTGBot/cmd/bot/config"
	"weatherTGBot/cmd/bot/providers"
	"weatherTGBot/pkg/db/postgres"
	"weatherTGBot/pkg/telegram"
)

const configPath = "cmd/bot/config/config.json"

func main() {
	ctx, _ := context.WithCancel(context.Background())
	cnf := config.LoadConfiguration(configPath)

	go func() {
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
		<-sigs
		fmt.Println("Exit and clean")
		ctx.Done()
	}()

	cleanup, err := initService(cnf)
	if err != nil {
		log.Fatalln(err)
	}

	<-ctx.Done()
	cleanup()
}

func initService(cnf config.Config) (func(), error) {
	//telegram
	bot, err := tgbotapi.NewBotAPI(cnf.Telegram.APIKey)
	if err != nil {
		return nil, err
	}
	bot.Debug = cnf.Telegram.Debug

	//weather API
	weather, err := owm.NewCurrent(cnf.Weather.Unit, cnf.Weather.Lang, cnf.Weather.APIKey)
	if err != nil {
		return nil, err
	}

	//Postgres
	db, err := postgres.NewDBConnection(config.GetPgDsn(cnf))
	if err != nil {
		return nil, err
		//log.Panic(err)
	}

	//logger
	logger := providers.ProvideConsoleLogger(cnf.Logger.Lvl)
	TelegramBot := telegram.NewBot(bot, weather, db, logger)

	if err = TelegramBot.Start(); err != nil {
		return nil, err
	}

	var cleaner = func() {
		fmt.Println("Start cleaning")
		TelegramBot.Stop()
		db.Close()
		//Если юзать логгер с файлом, то чистить сессию тоже
	}

	return cleaner, nil
}
