package main

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"os"
	"os/signal"
	"syscall"
	"weatherTGBot/cmd/bot/providers"
	"weatherTGBot/internal/config"
	"weatherTGBot/pkg/telegram"
)

const configPath = "internal/config/config.json"

func main() {
	ctx, _ := context.WithCancel(context.Background())
	cnf, err := config.LoadConfig(configPath)
	if err != nil {
		panic(err)
	}

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

func initService(cnf *config.Config) (func(), error) {
	//telegram
	//todo сделать провайдеров
	//logger
	logger := providers.ProvideConsoleLogger(cnf)

	bot, err := tgbotapi.NewBotAPI(cnf.TelegramApi.APIKey)
	if err != nil {
		return nil, err
	}
	bot.Debug = cnf.TelegramApi.Debug

	//weather API
	weatherApi, err := providers.ProvideWeatherApi(cnf)
	if err != nil {
		return nil, err
	}

	//TGBotRepository
	tgBotRepository, DBcloser, err := providers.ProvideTgBotRepo(cnf, logger)
	if err != nil {
		return nil, err
	}

	TelegramBot := telegram.NewBot(bot, weatherApi, tgBotRepository, logger)

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
