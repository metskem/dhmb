package misc

import (
	"fmt"
	"github.com/metskem/dhmb/db"
	"github.com/metskem/dhmb/exporter"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"sync"
	"time"
)

var mutex sync.Mutex

func Loop(m db.Monitor) {
	retries := 0
	matchPattern := regexp.MustCompile(m.ExpRespCode)
	for {
		transport := http.Transport{IdleConnTimeout: time.Second}
		client := http.Client{CheckRedirect: func(req *http.Request, via []*http.Request) error { return http.ErrUseLastResponse }, Timeout: time.Duration(m.Timeout) * time.Second, Transport: &transport}
		startTime := time.Now()
		resp, err := client.Get(m.Url)
		if resp != nil {
			_, _ = ioutil.ReadAll(resp.Body)
			resp.Body.Close()
		}
		elapsed := time.Since(startTime).Milliseconds()
		statusCode := 0
		errorString := ""
		if err == nil && resp != nil && matchPattern.MatchString(resp.Status) {
			log.Printf("%s: OK, statusCode: %d, respTime(ms): %d", m.MonName, resp.StatusCode, elapsed)
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
	mutex.Lock()
	defer mutex.Unlock()
	exporter.LastSeenMonitors[m.MonName] = exporter.LastSeenMonitor{Timestamp: time.Now(), Resptime: respTime, StatusUp: statusUp}
	monFromDB, err := db.GetMonitorByName(m.MonName)
	if err == nil {
		if monFromDB.LastStatus != db.MonLastStatusUp && statusUp {
			monFromDB.LastStatus = db.MonLastStatusUp
			monFromDB.LastStatusChanged = time.Now()
			_ = db.UpdateMonitor(monFromDB)
		}
		if monFromDB.LastStatus != db.MonLastStatusDown && !statusUp {
			monFromDB.LastStatus = db.MonLastStatusDown
			monFromDB.LastStatusChanged = time.Now()
			_ = db.UpdateMonitor(monFromDB)
		}
		recordResponseTime(m, respTime)
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
	if m.MonStatus != db.MonStatusSilenced {
		Broadcast(message)
	}
}

func recordResponseTime(m db.Monitor, respTime int64) {
	db.InsertRespTime(db.RespTime{MonId: m.Id, Timestamp: time.Now(), Time: respTime})
}
