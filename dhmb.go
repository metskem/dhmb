package main

import (
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/metskem/dhmb/conf"
	"github.com/metskem/dhmb/db"
	"log"
	"os"
	"strings"
	"time"
)

var me tgbotapi.User

func main() {
	token := os.Getenv("bottoken")
	if len(token) == 0 {
		log.Print("missing envvar \"bottoken\"")
		os.Exit(8)
	}

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err.Error())
	}

	if os.Getenv("debug") == "true" {
		bot.Debug = true
	}

	me, err = bot.GetMe()
	meDetails := "unknown"
	if err == nil {
		meDetails = fmt.Sprintf("BOT: ID:%d UserName:%s FirstName:%s LastName:%s", me.ID, me.UserName, me.FirstName, me.LastName)
		log.Printf("Started bot: %s, version:%s, build time:%s, commit hash:%s", meDetails, conf.VersionTag, conf.BuildTime, conf.CommitHash)
	} else {
		log.Printf("bot.GetMe() failed: %v", err)
	}

	database := db.Initdb()
	defer database.Close()

	newUpdate := tgbotapi.NewUpdate(0)
	newUpdate.Timeout = 60

	updatesChan, err := bot.GetUpdatesChan(newUpdate)
	if err == nil {

		// announce that we are alive again
		go func() {
			for i := 0; i < 1; i++ {
				time.Sleep(time.Second * 1)
				chatids := []int64{-235825137, 337345957} // TODO: for now fixed, but this should come from the chat table
				for _, chatid := range chatids {
					_, err := bot.Send(tgbotapi.NewMessage(chatid, fmt.Sprintf("%s started, buildtime: %s", meDetails, conf.BuildTime)))
					if err != nil {
						log.Printf("failed sending message to chat %d, error is %v", chatid, err)
					}
				}
			}
		}()

		// start listening for messages, and optionally respond
		for update := range updatesChan {
			if update.Message == nil { // ignore any non-Message Updates
				log.Println("ignored null update")
			} else {
				chat := update.Message.Chat
				mentionedMe, cmdMe := talkOrCmdToMe(update)
				if chat.IsPrivate() || (chat.IsGroup() && mentionedMe) {
					log.Printf("[%s] [chat:%d] %s\n", update.Message.From.UserName, chat.ID, update.Message.Text)
					if cmdMe {

						// do the actual send Message
						_, err := bot.Send(tgbotapi.NewMessage(chat.ID, fmt.Sprintf("Hi user %s, your name is %s %s", update.Message.ReplyToMessage.Chat.UserName, update.Message.ReplyToMessage.Chat.FirstName, update.Message.ReplyToMessage.Chat.LastName)))

						if err != nil {
							log.Printf("failed sending message: %v", err)
						}
					}
				}
			}
			fmt.Println("")
		}
	} else {
		log.Printf("failed getting bot updatesChannel, error: %v", err)
		os.Exit(8)
	}
}

/*
  Returns if we are mentioned and if we were commanded
*/
func talkOrCmdToMe(update tgbotapi.Update) (bool, bool) {
	entities := update.Message.Entities
	var mentioned = false
	var botCmd = false
	if entities != nil {
		for _, entity := range *entities {
			if entity.Type == "mention" {
				if strings.HasPrefix(update.Message.Text, fmt.Sprintf("@%s", me.UserName)) {
					mentioned = true
				}
			}
			if entity.Type == "bot_command" {
				botCmd = true
			}
		}
	}
	// if another bot was mentioned, the cmd is not for us
	if update.Message.Chat.IsGroup() && mentioned == false {
		botCmd = false
	}
	return mentioned, botCmd
}
