package main

import (
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api"
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
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}

	if os.Getenv("debug") != "" {
		bot.Debug = true
	}

	me, err := bot.GetMe()
	if err == nil {
		log.Printf("I am bot: ID:%d UserName:%s FirstName:%s LastName:%s IsBot:%t LanguageCode:%s \n", me.ID, me.UserName, me.FirstName, me.LastName, me.IsBot, me.LanguageCode)
	}

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	go func() {
		for i := 0; i < 5; i++ {
			time.Sleep(time.Duration(time.Second * 5))
			bot.Send(tgbotapi.NewMessage(-235825137, fmt.Sprintf("it is now %v", time.Now())))
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
