package main

import (
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/metskem/dhmb/check"
	"github.com/metskem/dhmb/conf"
	"github.com/metskem/dhmb/db"
	"log"
	"os"
	"strings"
)

var me tgbotapi.User

func main() {
	token := os.Getenv("bottoken")
	if len(token) == 0 {
		log.Print("missing envvar \"bottoken\"")
		os.Exit(8)
	}

	var err error

	conf.Bot, err = tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err.Error())
	}

	if os.Getenv("debug") == "true" {
		conf.Bot.Debug = true
	}

	me, err = conf.Bot.GetMe()
	meDetails := "unknown"
	if err == nil {
		meDetails = fmt.Sprintf("BOT: ID:%d UserName:%s FirstName:%s LastName:%s", me.ID, me.UserName, me.FirstName, me.LastName)
		log.Printf("Started Bot: %s, version:%s, build time:%s, commit hash:%s", meDetails, conf.VersionTag, conf.BuildTime, conf.CommitHash)
	} else {
		log.Printf("Bot.GetMe() failed: %v", err)
	}

	db.Initdb()

	newUpdate := tgbotapi.NewUpdate(0)
	newUpdate.Timeout = 60

	updatesChan, err := conf.Bot.GetUpdatesChan(newUpdate)
	if err == nil {

		// announce that we are live again
		chats := db.GetChats()
		for _, chat := range chats {
			_, err := conf.Bot.Send(tgbotapi.NewMessage(chat.ChatId, fmt.Sprintf("%s has been restarted, buildtime: %s", me.UserName, conf.BuildTime)))
			if err != nil {
				log.Printf("failed sending message to chat %d, error is %v", chat.ChatId, err)
			}
		}

		// start the checks
		check.Runner()

		// start listening for messages, and optionally respond
		for update := range updatesChan {
			if update.Message == nil { // ignore any non-Message Updates
				log.Println("ignored null update")
			} else {
				chat := update.Message.Chat
				mentionedMe, cmdMe := talkOrCmdToMe(update)

				// check if someone is talking to me:
				if chat.IsPrivate() || (chat.IsGroup() && mentionedMe) {
					log.Printf("[%s] [chat:%d] %s\n", update.Message.From.UserName, chat.ID, update.Message.Text)
					if cmdMe {
						// do the actual send Message
						_, err := conf.Bot.Send(tgbotapi.NewMessage(chat.ID, fmt.Sprintf("Hi user %s, your name is %s %s", update.Message.From.UserName, update.Message.From.FirstName, update.Message.From.LastName)))
						if err != nil {
							log.Printf("failed sending message: %v", err)
						}
					}
				}

				// check if someone started a new chat or added me to a group
				if update.Message.NewChatMembers != nil && len(*update.Message.NewChatMembers) > 0 {
					for _, user := range *update.Message.NewChatMembers {
						if user.UserName == me.UserName {
							name := chat.UserName
							if chat.IsGroup() {
								name = chat.Title
							}
							log.Printf("new chat added, chatid: %d, chat UserName: %s\n", chat.ID, name)
						}
					}
				}

				// check if someone deleted a chat or removed me from a group
				if update.Message.LeftChatMember != nil {
					leftChatMember := *update.Message.LeftChatMember
					if leftChatMember.UserName == me.UserName {
						name := chat.UserName
						if chat.IsGroup() {
							name = chat.Title
						}
						log.Printf("chat removed, chatid: %d, chat UserName: %s\n", chat.ID, name)
					}
				}
			}
			fmt.Println("")
		}
	} else {
		log.Printf("failed getting Bot updatesChannel, error: %v", err)
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
