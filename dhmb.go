package main

import (
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/metskem/dhmb/conf"
	"github.com/metskem/dhmb/db"
	"log"
	"os"
	"time"
)

func main() {
	token := os.Getenv("bottoken")
	if len(token) == 0 {
		log.Print("missing envvar \"bottoken\"")
		os.Exit(8)
	}

	database := db.Initdb()
	defer database.Close()

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}

	if os.Getenv("debug") == "true" {
		bot.Debug = true
	}

	me, err := bot.GetMe()
	meDetails := fmt.Sprintf("BOT: ID:%d UserName:%s FirstName:%s LastName:%s", me.ID, me.UserName, me.FirstName, me.LastName)
	if err == nil {
		log.Printf("Started bot: %s, version:%s, build time:%s, commit hash:%s", meDetails, conf.VersionTag, conf.BuildTime, conf.CommitHash)
	}

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	go func() {
		for i := 0; i < 1; i++ {
			time.Sleep(time.Duration(time.Second * 5))
			bot.Send(tgbotapi.NewMessage(-235825137, fmt.Sprintf("%s started", meDetails)))
			bot.Send(tgbotapi.NewMessage(1140134411, fmt.Sprintf("%s started", meDetails)))
		}
	}()

	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			log.Println("ignored null update")
		}

		log.Printf("[%s] %s\n", update.Message.From.UserName, update.Message.Text)

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
		//msg.ReplyToMessageID = update.Message.MessageID

		bot.Send(msg)
	}
}
