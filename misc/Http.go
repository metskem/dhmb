package misc

import (
	"fmt"
	"github.com/metskem/dhmb/db"
	"github.com/metskem/dhmb/exporter"
	"io"
	"log"
	"math/rand"
	"net/http"
	"regexp"
	"sync"
	"time"
)

var mutex sync.Mutex

func Loop(m db.Monitor) {
	retries := 0
	matchPatternResponse := regexp.MustCompile(m.ExpResponse)
	matchPatternRespCode := regexp.MustCompile(m.ExpRespCode)
	for {
		transport := http.Transport{IdleConnTimeout: time.Second, DisableKeepAlives: true}
		client := http.Client{CheckRedirect: func(req *http.Request, via []*http.Request) error { return http.ErrUseLastResponse }, Timeout: time.Duration(m.Timeout) * time.Second, Transport: &transport}
		req := setRandomUseragent()
		req.URL, _ = req.URL.Parse(m.Url)
		startTime := time.Now()
		resp, err := client.Do(req)
		elapsed := time.Since(startTime).Milliseconds()
		statusCode := 0
		errorString := ""
		if err == nil && resp != nil && matchPatternRespCode.MatchString(resp.Status) {
			respBody, _ := io.ReadAll(resp.Body)
			if matchPatternResponse.MatchString(string(respBody)) {
				log.Printf("%s: OK, statusCode: %d, respTime(ms): %d", m.MonName, resp.StatusCode, elapsed)
				updateLastStatus(m, true, elapsed)
				if retries >= m.Retries {
					alert(true, m, statusCode, "")
				}
				retries = 0
			} else {
				log.Printf("%s: NOK, pattern \"%s\" not found in response body, statusCode: %d, respTime(ms): %d", m.MonName, m.ExpResponse, resp.StatusCode, elapsed)
				updateLastStatus(m, false, elapsed)
				retries++
				if retries == m.Retries {
					alert(false, m, statusCode, errorString)
				}
			}
			_ = resp.Body.Close()
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

func setRandomUseragent() (req *http.Request) {
	req, err := http.NewRequest("GET", "https://example.com", nil)
	if err != nil {
		log.Printf("Error creating request: %s", err)
		return
	}
	req.Header.Set("User-Agent", fmt.Sprintf("dhmb-client/1.%d", rand.Intn(100)))
	return req
}

func updateLastStatus(m db.Monitor, statusUp bool, respTime int64) {
	mutex.Lock()
	defer mutex.Unlock()
	exporter.LastSeenMonitors[m.MonName] = exporter.LastSeenMonitor{Timestamp: time.Now(), RespTime: respTime, StatusUp: statusUp}
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
