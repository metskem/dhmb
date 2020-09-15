package check

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/metskem/dhmb/conf"
	"github.com/metskem/dhmb/db"
	"log"
	"net/http"
	"time"
)

func Loop(m db.Monitor) {
	retries := 0
	for m.MonStatus == db.MonStatusActive {
		client := http.Client{Timeout: time.Duration(m.Timeout) * time.Second}
		resp, err := client.Get(m.Url)
		statusCode := 0
		errorString := ""
		if err == nil && resp != nil && resp.StatusCode == m.ExpRespCode {
			log.Printf("%s:  OK, resp code: %d", m, resp.StatusCode)
			retries = 0
		} else {
			if resp != nil {
				statusCode = resp.StatusCode
				errorString = resp.Status
			}
			if err != nil {
				errorString = err.Error()
			}
			log.Printf("%s: NOK, resp code: %d, error getting URL %s: %s", m, statusCode, m.Url, errorString)
			retries++
			if retries == m.Retries {
				alert(m, statusCode, errorString)
			}
		}
		time.Sleep(time.Duration(m.Interval) * time.Second)
	}
}

func alert(m db.Monitor, statusCode int, errorString string) {
	for _, chat := range db.GetChats() {
		_, err := conf.Bot.Send(tgbotapi.NewMessage(chat.ChatId, fmt.Sprintf("%s is down, statusCode: %d, error: %s", m.MonName, statusCode, errorString)))
		if err != nil {
			log.Printf("failed sending message to chat %d, error is %v", chat.ChatId, err)
		}

	}
}
