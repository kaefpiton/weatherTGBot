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
	//todo сделать провайдеров
	bot, err := tgbotapi.NewBotAPI(cnf.TelegramApi.APIKey)
	if err != nil {
		return nil, err
	}
	bot.Debug = cnf.TelegramApi.Debug

	//weather API
	weather, err := owm.NewCurrent(cnf.WeatherApi.Unit, cnf.WeatherApi.Lang, cnf.WeatherApi.APIKey)
	if err != nil {
		return nil, err
	}

	//TGBotRepository
	tgBotRepository, DBcloser, err := providers.ProvideTgBotRepo(cnf)
	if err != nil {
		return nil, err
	}

	//logger
	logger := providers.ProvideConsoleLogger(cnf)
	TelegramBot := telegram.NewBot(bot, weather, tgBotRepository, logger)

	if err = TelegramBot.Start(); err != nil {
		return nil, err
	}

	var cleaner = func() {
		fmt.Println("Start cleaning")
		TelegramBot.Stop()
		DBcloser()
		//Если юзать логгер с файлом, то закрывать сессию тоже
	}

	return cleaner, nil
}
