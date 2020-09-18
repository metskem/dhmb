package misc

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/metskem/dhmb/db"
	"log"
	"strings"
)

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
		msg := "restart function not implemented yet"
		log.Println(msg)
		SendMessage(db.Chat{ChatId: update.Message.Chat.ID}, msg)
	}
	if strings.HasPrefix(update.Message.Text, "/status") {
		msg := "status function not implemented yet"
		log.Println(msg)
		SendMessage(db.Chat{ChatId: update.Message.Chat.ID}, msg)
	}
	if strings.HasPrefix(update.Message.Text, "/start") {
		SendMessage(db.Chat{ChatId: update.Message.Chat.ID}, fmt.Sprintf("Hi %s (%s %s), you will receive alerts from now", update.Message.From.UserName, update.Message.From.FirstName, update.Message.From.LastName))
	}
}
