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

var Me tgbotapi.User
var NumRunningMonitors int
var RestartRequested = false
var MonCountLock = sync.RWMutex{}

var Bot *tgbotapi.BotAPI

func SendMessage(chat db.Chat, message string) {
	_, err := Bot.Send(tgbotapi.NewMessage(chat.ChatId, message))
	if err != nil {
		log.Printf("failed sending message to chat %d, error is %v", chat.ChatId, err)
		if err.Error() == "Forbidden: bot was blocked by the user" || err.Error() == "Forbidden: bot was kicked from the group chat" {
			if db.DeleteChat(chat.ChatId) {
				log.Printf("deleted chatid %d from list", chat.ChatId)
				Broadcast(fmt.Sprintf("chat removed: (Id=%d)", chat.ChatId))
			} else {
				log.Printf("not deleted chatid %d from list", chat.ChatId)
			}
		}
	}
}

func HandleCommand(update tgbotapi.Update) {
	chatter := db.Chat{ChatId: update.Message.Chat.ID}
	if strings.HasPrefix(update.Message.Text, "/restart") {
		msg := fmt.Sprintf("restart requested by %s, please wait...", update.Message.From.UserName)
		log.Println(msg)
		SendMessage(chatter, msg)
		RestartRequested = true
		Runner()
		log.Println("restart completed")
		Broadcast(fmt.Sprintf("restart done by %s", update.Message.From.UserName))
	}

	if strings.HasPrefix(update.Message.Text, "/status") {
		var msg string
		for ix, mon := range db.GetActiveMonitors() {
			msg = fmt.Sprintf("%s%d - %s: %s since %s\n", msg, ix, mon.MonName, mon.LastStatus, mon.LastStatusChanged.Format(time.RFC3339))
		}
		log.Println("\n" + msg)
		SendMessage(chatter, msg)
	}

	if strings.HasPrefix(update.Message.Text, "/members") {
		var msg string
		for ix, dhmbChat := range db.GetChats() {
			chat, err := Bot.GetChat(tgbotapi.ChatConfig{ChatID: dhmbChat.ChatId})
			if err == nil {
				if chat.IsGroup() {
					msg = fmt.Sprintf("%s%d - chat:%d,  group: %s (%s)\n", msg, ix, chat.ID, chat.Title, chat.Description)
				} else {
					msg = fmt.Sprintf("%s%d - chat: %d,  user : %s (%s %s)\n", msg, ix, chat.ID, chat.UserName, chat.FirstName, chat.LastName)
				}
			} else {
				log.Printf("error getting chat %d: %v", dhmbChat.ChatId, err)
			}
		}
		log.Println("\n" + msg)
		SendMessage(chatter, msg)
	}

	if strings.HasPrefix(update.Message.Text, "/start") {
		SendMessage(chatter, fmt.Sprintf("Hi %s (%s %s), you will receive alerts from now", update.Message.From.UserName, update.Message.From.FirstName, update.Message.From.LastName))
	}

	if strings.HasPrefix(update.Message.Text, "/debug") {
		if strings.Contains(update.Message.Text, " on") {
			Bot.Debug = true
			SendMessage(chatter, "debug turned on")
		} else {
			if strings.Contains(update.Message.Text, " off") {
				Bot.Debug = false
				SendMessage(chatter, "debug turned off")
			} else {
				SendMessage(chatter, "please specify /debug on  or  /debug off")
			}
		}
	}

	if strings.HasPrefix(update.Message.Text, "/usernames") {
		var msg string
		for ix, userName := range db.GetUserNames() {
			msg = fmt.Sprintf("%s%d - %s  :  %s\n", msg, ix, userName.Name, userName.Role)
		}
		log.Println("\n" + msg)
		SendMessage(chatter, msg)
	}

	if strings.HasPrefix(update.Message.Text, "/usernameadd") {
		words := strings.Split(update.Message.Text, " ")
		if len(words) == 3 && (words[2] == db.UserNameRoleAdmin || words[2] == db.UserNameRoleReader) {
			if strings.HasSuffix(update.Message.Text, fmt.Sprintf(" %s", db.UserNameRoleReader)) {
				if db.InsertUserName(db.UserName{Name: words[1], Role: db.UserNameRoleReader}) != 0 {
					SendMessage(chatter, fmt.Sprintf("username %s with role %s added", words[1], words[2]))
					log.Printf("username %s with role %s added", words[1], words[2])
				} else {
					SendMessage(chatter, fmt.Sprintf("username %s with role %s not added", words[1], words[2]))
					log.Printf("username %s with role %s not added", words[1], words[2])
				}
				return
			} else {
				if strings.HasSuffix(update.Message.Text, fmt.Sprintf(" %s", db.UserNameRoleAdmin)) {
					if db.InsertUserName(db.UserName{Name: words[1], Role: db.UserNameRoleAdmin}) != 0 {
						SendMessage(chatter, fmt.Sprintf("username %s with role %s added", words[1], words[2]))
						log.Printf("username %s with role %s added", words[1], words[2])
					} else {
						SendMessage(chatter, fmt.Sprintf("username %s with role %s not added", words[1], words[2]))
						log.Printf("username %s with role %s not added", words[1], words[2])
					}
					return
				}
			}
		}
		SendMessage(chatter, "specify /usernameadd <username> [admin|reader]")
	}

	if strings.HasPrefix(update.Message.Text, "/usernamedelete") {
		words := strings.Split(update.Message.Text, " ")
		if len(words) == 2 {
			if db.DeleteUserName(words[1]) {
				SendMessage(chatter, fmt.Sprintf("username %s deleted", words[1]))
				log.Printf("username %s deleted", words[1])
			} else {
				SendMessage(chatter, fmt.Sprintf("username %s not deleted", words[1]))
				log.Printf("username %s not deleted", words[1])
			}
		} else {
			SendMessage(chatter, "specify /usernamedelete <username>")
		}
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

func HasRole(userName string, roleName string) bool {
	for _, dbuser := range db.GetUserNames() {
		if dbuser.Name == userName && dbuser.Role == roleName {
			return true
		}
	}
	log.Printf("%s permission denied for user %s", roleName, userName)
	return false
}

/*
  Returns if we are mentioned and if we were commanded
*/
func TalkOrCmdToMe(update tgbotapi.Update) (bool, bool) {
	entities := update.Message.Entities
	var mentioned = false
	var botCmd = false
	if entities != nil {
		for _, entity := range *entities {
			if entity.Type == "mention" {
				if strings.HasPrefix(update.Message.Text, fmt.Sprintf("@%s", Me.UserName)) {
					mentioned = true
				}
			}
			if entity.Type == "bot_command" {
				botCmd = true
				if strings.Contains(update.Message.Text, fmt.Sprintf("@%s", Me.UserName)) {
					mentioned = true
				}
			}
		}
	}
	// if another bot was mentioned, the cmd is not for us
	if update.Message.Chat.IsGroup() && mentioned == false {
		botCmd = false
	}
	return mentioned, botCmd
}

func Broadcast(message string) {
	for _, chat := range db.GetChats() {
		if db.UserNameIsAdmin(chat.Name) {
			SendMessage(chat, message)
		}
	}
}
