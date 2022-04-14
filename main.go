package main

import (
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/metskem/dhmb/conf"
	"github.com/metskem/dhmb/db"
	"github.com/metskem/dhmb/exporter"
	"github.com/metskem/dhmb/misc"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

func main() {

	conf.EnvironmentComplete()
	log.SetOutput(os.Stdout)

	var err error

	misc.Bot, err = tgbotapi.NewBotAPI(conf.BotToken)
	if err != nil {
		log.Panic(err.Error())
	}

	misc.Bot.Debug = conf.Debug

	misc.Me, err = misc.Bot.GetMe()
	meDetails := "unknown"
	if err == nil {
		meDetails = fmt.Sprintf("BOT: ID:%d UserName:%s FirstName:%s LastName:%s", misc.Me.ID, misc.Me.UserName, misc.Me.FirstName, misc.Me.LastName)
		log.Printf("Started Bot: %s, version:%s, build time:%s, commit hash:%s", meDetails, conf.VersionTag, conf.BuildTime, conf.CommitHash)
	} else {
		log.Printf("Bot.GetMe() failed: %v", err)
	}

	db.Initdb()

	// fire up a background thread that regularly deletes old rows from the resptime table
	go func() {
		for {
			time.Sleep(time.Minute * 10)
			monitors, err := db.GetMonitorsByStatus(db.MonStatusActive)
			if err == nil {
				for _, mon := range monitors {
					db.CleanupOldStuffForMonitor(mon)
				}
			} else {
				log.Printf("failed to get monitors: %s", err)
			}
		}
	}()

	// fire up a background thread that monitors the currently running routines by checking if the last resptime update is recent enough
	go func() {
		for {
			time.Sleep(time.Minute * 12)
			misc.CheckLastRespTimeUpdates()
		}
	}()

	// start the prometheus exporter
	go func() {
		//Create a new instance of the collector and register with the exporter client.
		collector := exporter.NewDHMBbCollector()
		prometheus.MustRegister(collector)

		// expose metrics on the /metrics endpoint.
		http.Handle("/metrics", promhttp.Handler())
		log.Printf("Prometheus exporter on port %d", conf.PromExporterPort)
		log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", conf.PromExporterPort), nil))
	}()

	newUpdate := tgbotapi.NewUpdate(0)
	newUpdate.Timeout = 60

	updatesChan, err := misc.Bot.GetUpdatesChan(newUpdate)
	if err == nil {

		// announce that we are live again
		misc.Broadcast(fmt.Sprintf("%s has been (re)started, buildtime: %s", misc.Me.UserName, conf.BuildTime))

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
						// /status can be done by anyone, for the other cmds you need admin role
						if misc.HasRole(fromUser, db.UserNameRoleAdmin) || strings.HasPrefix(update.Message.Text, "/status") || strings.HasPrefix(update.Message.Text, "/chart") {
							misc.HandleCommand(update)
						} else {
							misc.SendMessage(db.Chat{ChatId: chat.ID}, fmt.Sprintf("sorry, %s is not allowed to send me that command", fromUser))
						}
					}
				}

				// check if someone started a new chat
				if chat.IsPrivate() && cmdMe && update.Message.Text == "/start" {
					if db.InsertChat(db.Chat{ChatId: chat.ID, Name: chat.UserName}) != 0 {
						log.Printf("new chat added, chatid: %d, chat: %s (%s %s)\n", chat.ID, chat.UserName, chat.FirstName, chat.LastName)
						misc.Broadcast(fmt.Sprintf("new member: chat: %s (%s %s)", chat.UserName, chat.FirstName, chat.LastName))
					}
				}

				// check if someone added me to a group
				if update.Message.NewChatMembers != nil && len(*update.Message.NewChatMembers) > 0 {
					if misc.HasRole(update.Message.From.UserName, db.UserNameRoleReader) || misc.HasRole(update.Message.From.UserName, db.UserNameRoleAdmin) {
						for _, user := range *update.Message.NewChatMembers {
							if user.UserName == misc.Me.UserName {
								if db.InsertChat(db.Chat{ChatId: chat.ID, Name: chat.UserName}) != 0 {
									log.Printf("new chat added, chatid: %d, chat: %s (%s %s)\n", chat.ID, chat.Title, chat.FirstName, chat.LastName)
									misc.Broadcast(fmt.Sprintf("new member: group:%s", chat.Title))
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
							if db.DeleteChat(chat.ID) {
								log.Printf("chat removed, chatid: %d, chat: %s (%s %s)\n", chat.ID, chat.Title, chat.FirstName, chat.LastName)
								misc.Broadcast(fmt.Sprintf("chat removed: %s (%s %s)", chat.UserName, chat.FirstName, chat.LastName))
							} else {
								log.Printf("chat not deleted, chatid: %d, chat: %s (%s %s)\n", chat.ID, chat.Title, chat.FirstName, chat.LastName)
							}
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
