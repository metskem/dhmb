package check

import (
	"fmt"
	"github.com/metskem/dhmb/db"
	"github.com/metskem/dhmb/misc"
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
			log.Printf("%s: OK, statusCode: %d", m.MonName, resp.StatusCode)
			if retries >= m.Retries {
				alert(true, m, statusCode, "")
			}
			retries = 0
		} else {
			if resp != nil {
				statusCode = resp.StatusCode
				errorString = resp.Status
			}
			if err != nil {
				errorString = err.Error()
			}
			log.Printf("%s: NOK, attempt %d, statusCode: %d, error getting URL %s: %s", m.MonName, retries, statusCode, m.Url, errorString)
			retries++
			if retries == m.Retries {
				alert(false, m, statusCode, errorString)
			}
		}
		time.Sleep(time.Duration(m.Interval) * time.Second)
	}
}

func alert(statusUp bool, m db.Monitor, statusCode int, errorString string) {
	var message string
	if statusUp {
		message = fmt.Sprintf("%s is UP: statusCode: %d", m.MonName, statusCode)
	} else {
		message = fmt.Sprintf("%s is DOWN: statusCode: %d, error: %s", m.MonName, statusCode, errorString)
	}
	log.Println(message)
	for _, chat := range db.GetChats() {
		misc.SendMessage(chat, message)
	}
}
