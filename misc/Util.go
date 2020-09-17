package misc

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/metskem/dhmb/db"
	"log"
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
