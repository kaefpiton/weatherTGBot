package main

import (
	"fmt"
	owm "github.com/briandowns/openweathermap"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"weatherTGBot/cmd/bot/config"
	"weatherTGBot/pkg/db/postgres"
	"weatherTGBot/pkg/telegram"
)


func main(){
	config := config.LoadConfiguration("config/config.json")

	bot, err := tgbotapi.NewBotAPI(config.TelegramAPIKey)
	if err != nil {
		log.Panic(err)
	}
	bot.Debug = true


	weather, err := owm.NewCurrent(config.Weather.Unit, config.Weather.Lang, config.Weather.APIKey)
	if err != nil {
		log.Fatalln(err)
	}


	dataSource :=fmt.Sprintf("user=%s password=%s dbname=%s sslmode=%s",config.DataBase.User,
																config.DataBase.Password,
																config.DataBase.DBName,
																config.DataBase.SSLMode)
	db,err:= postgres.NewDB(dataSource)
	if err != nil {
		log.Panic(err)
	}

	TelegramBot := telegram.NewBot(bot, weather, db)

	if err := TelegramBot.Start(); err != nil{
		log.Fatal(err)
	}

}

