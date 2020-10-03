package main

import (
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/metskem/dhmb/conf"
	"github.com/metskem/dhmb/db"
	"github.com/metskem/dhmb/misc"
	"log"
	"os"
)

func main() {
	token := os.Getenv("bottoken")
	if len(token) == 0 {
		log.Print("missing envvar \"bottoken\"")
		os.Exit(8)
	}

	var err error

	misc.Bot, err = tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err.Error())
	}

	if os.Getenv("debug") == "true" {
		misc.Bot.Debug = true
	}

	misc.Me, err = misc.Bot.GetMe()
	meDetails := "unknown"
	if err == nil {
		meDetails = fmt.Sprintf("BOT: ID:%d UserName:%s FirstName:%s LastName:%s", misc.Me.ID, misc.Me.UserName, misc.Me.FirstName, misc.Me.LastName)
		log.Printf("Started Bot: %s, version:%s, build time:%s, commit hash:%s", meDetails, conf.VersionTag, conf.BuildTime, conf.CommitHash)
	} else {
		log.Printf("Bot.GetMe() failed: %v", err)
	}

	db.Initdb()

	newUpdate := tgbotapi.NewUpdate(0)
	newUpdate.Timeout = 60

	updatesChan, err := misc.Bot.GetUpdatesChan(newUpdate)
	if err == nil {

		// announce that we are live again
		chats := db.GetChats()
		for _, chat := range chats {
			misc.SendMessage(chat, fmt.Sprintf("%s has been restarted, buildtime: %s", misc.Me.UserName, conf.BuildTime))
		}

		// start the checks
		misc.Runner()

		// start listening for messages, and optionally respond
		for update := range updatesChan {
			if update.Message == nil { // ignore any non-Message Updates
				log.Println("ignored null update")
			} else {
				chat := update.Message.Chat
				mentionedMe, cmdMe := misc.TalkOrCmdToMe(update)

				// check if someone is talking to me:
				if (chat.IsPrivate() || (chat.IsGroup() && mentionedMe)) && update.Message.Text != "/start" {
					log.Printf("[%s] [chat:%d] %s\n", update.Message.From.UserName, chat.ID, update.Message.Text)
					if cmdMe {
						fromUser := update.Message.From.UserName
						if chat.IsPrivate() {
							fromUser = chat.UserName
						}
						if misc.HasRole(fromUser, db.UserNameRoleAdmin) {
							misc.HandleCommand(update)
						} else {
							misc.SendMessage(db.Chat{ChatId: chat.ID}, fmt.Sprintf("sorry, %s is not allowed to send me commands", fromUser))
						}
					}
				}

				// check if someone started a new chat
				if chat.IsPrivate() && cmdMe && update.Message.Text == "/start" {
					if misc.HasRole(chat.UserName, db.UserNameRoleReader) || misc.HasRole(chat.UserName, db.UserNameRoleAdmin) {
						if db.InsertChat(db.Chat{ChatId: chat.ID}) != 0 {
							log.Printf("new chat added, chatid: %d, chat: %s (%s %s)\n", chat.ID, chat.UserName, chat.FirstName, chat.LastName)
						}
					} else {
						misc.SendMessage(db.Chat{ChatId: chat.ID}, "sorry, you are not allowed to talk or listen to me")
					}
				}

				// check if someone added me to a group
				if update.Message.NewChatMembers != nil && len(*update.Message.NewChatMembers) > 0 {
					if misc.HasRole(update.Message.From.UserName, db.UserNameRoleReader) {
						for _, user := range *update.Message.NewChatMembers {
							if user.UserName == misc.Me.UserName {
								if db.InsertChat(db.Chat{ChatId: chat.ID}) != 0 {
									log.Printf("new chat added, chatid: %d, chat: %s (%s %s)\n", chat.ID, chat.Title, chat.FirstName, chat.LastName)
								}
							}
						}
					} else {
						misc.SendMessage(db.Chat{ChatId: chat.ID}, fmt.Sprintf("sorry, %s is not allowed to add me to a group", update.Message.From.UserName))
					}
				}

				// check if someone removed me from a group
				if update.Message.LeftChatMember != nil {
					if misc.HasRole(update.Message.From.UserName, db.UserNameRoleReader) {
						leftChatMember := *update.Message.LeftChatMember
						if leftChatMember.UserName == misc.Me.UserName {
							db.DeleteChat(chat.ID)
							log.Printf("chat removed, chatid: %d, chat: %s (%s %s)\n", chat.ID, chat.Title, chat.FirstName, chat.LastName)
						}
					} else {
						misc.SendMessage(db.Chat{ChatId: chat.ID}, fmt.Sprintf("sorry, %s is not allowed to remove me from a group", update.Message.From.UserName))
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
