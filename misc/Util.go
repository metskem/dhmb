package misc

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/metskem/dhmb/db"
	"log"
	"strings"
	"sync"
	"time"
)

var NumRunningMonitors int
var RestartRequested = false
var MonCountLock = sync.RWMutex{}

var Bot *tgbotapi.BotAPI

func SendMessage(chat db.Chat, message string) {
	_, err := Bot.Send(tgbotapi.NewMessage(chat.ChatId, message))
	if err != nil {
		log.Printf("failed sending message to chat %d, error is %v", chat.ChatId, err)
		if err.Error() == "Forbidden: bot was blocked by the user" {
			db.DeleteChat(chat.ChatId)
			log.Printf("removed chatid %d from list", chat.ChatId)
		}
	}
}

func HandleCommand(update tgbotapi.Update) {

	if strings.HasPrefix(update.Message.Text, "/restart") {
		msg := "restart function accepted, please wait..."
		log.Println(msg)
		SendMessage(db.Chat{ChatId: update.Message.Chat.ID}, msg)
		RestartRequested = true
		Runner()
		msg = "restart function completed"
		log.Println(msg)
		SendMessage(db.Chat{ChatId: update.Message.Chat.ID}, msg)
	}

	if strings.HasPrefix(update.Message.Text, "/status") {
		var msg string
		for ix, mon := range db.GetActiveMonitors() {
			msg = fmt.Sprintf("%s%d - %s: %s since %s\n", msg, ix, mon.MonName, mon.LastStatus, mon.LastStatusChanged.Format(time.RFC3339))
		}
		log.Println("\n" + msg)
		SendMessage(db.Chat{ChatId: update.Message.Chat.ID}, msg)
	}

	if strings.HasPrefix(update.Message.Text, "/start") {
		SendMessage(db.Chat{ChatId: update.Message.Chat.ID}, fmt.Sprintf("Hi %s (%s %s), you will receive alerts from now", update.Message.From.UserName, update.Message.From.FirstName, update.Message.From.LastName))
	}
}

/**
Wait for the monitor interval to expire before returning with false. If returned with true, a restart is requested, and the caller can decide to no longer loop.
*/
func RestartOrWait(m db.Monitor) bool {
	endWaitTime := time.Now().Add(time.Duration(m.Interval) * time.Second)
	for true {
		if RestartRequested {
			MonCountLock.Lock()
			NumRunningMonitors--
			log.Printf("stopped %s, num monitors left is %d", m, NumRunningMonitors)
			MonCountLock.Unlock()
			return true
		}
		if time.Now().Before(endWaitTime) {
			time.Sleep(time.Second * 3)
		} else {
			return false
		}
	}
	return false
}

/*
  Iterate over the monitors in the monitor table, and start a separate go routine for a check
*/
func Runner() {
	for RestartRequested && NumRunningMonitors != 0 {
		time.Sleep(time.Second * 3)
		log.Printf("waiting for restart to complete, number of running monitors is %d...\n", NumRunningMonitors)
	}
	RestartRequested = false
	for _, m := range db.GetActiveMonitors() {
		monitor := m
		go func(db.Monitor) {
			if monitor.MonType == db.MonTypeHttp {
				Loop(monitor)
			}
		}(monitor)
		MonCountLock.Lock()
		NumRunningMonitors++
		log.Printf("started %s, num monitors is %d", m, NumRunningMonitors)
		MonCountLock.Unlock()
	}
	log.Printf("we have %d running monitors", NumRunningMonitors)
}
