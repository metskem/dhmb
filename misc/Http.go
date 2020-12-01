package misc

import (
	"fmt"
	"github.com/metskem/dhmb/db"
	"log"
	"net/http"
	"time"
)

func Loop(m db.Monitor) {
	retries := 0
	for {
		client := http.Client{Timeout: time.Duration(m.Timeout) * time.Second}
		startTime := time.Now()
		resp, err := client.Get(m.Url)
		elapsed := int64(time.Since(startTime))
		statusCode := 0
		errorString := ""
		if err == nil && resp != nil && resp.StatusCode == m.ExpRespCode {
			log.Printf("%s: OK, statusCode: %d", m.MonName, resp.StatusCode)
			updateLastStatus(m, true, elapsed)
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
			updateLastStatus(m, false, elapsed)
			retries++
			if retries == m.Retries {
				alert(false, m, statusCode, errorString)
			}
		}
		if RestartOrWait(m) {
			return
		}
	}
}

func updateLastStatus(m db.Monitor, statusUp bool, respTime int64) {
	monFromDB := db.GetMonitorByName(m.MonName)
	if monFromDB.LastStatus != db.MonLastStatusUp && statusUp {
		monFromDB.LastStatus = db.MonLastStatusUp
		monFromDB.LastStatusChanged = time.Now()
		db.UpdateMonitor(monFromDB)
	}
	if monFromDB.LastStatus != db.MonLastStatusDown && !statusUp {
		monFromDB.LastStatus = db.MonLastStatusDown
		monFromDB.LastStatusChanged = time.Now()
		db.UpdateMonitor(monFromDB)
	}
	recordResponseTime(m, respTime)
}

func alert(statusUp bool, m db.Monitor, statusCode int, errorString string) {
	var message string
	if statusUp {
		message = fmt.Sprintf("%s is UP: statusCode: %d", m.MonName, statusCode)
	} else {
		message = fmt.Sprintf("%s is DOWN: statusCode: %d, error: %s", m.MonName, statusCode, errorString)
	}
	log.Println(message)
	if m.MonStatus != db.MonStatusSilenced {
		Broadcast(message)
	}
}

func recordResponseTime(m db.Monitor, respTime int64) {
	db.InsertRespTime(db.RespTime{MonId: m.Id, Time: respTime})
}
