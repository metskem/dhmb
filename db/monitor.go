package db

import (
	"errors"
	"fmt"
	"log"
	"time"
)

type Monitor struct {
	Id                int
	MonName           string
	MonType           string
	MonStatus         string
	Url               string
	Interval          int
	ExpRespCode       string
	Timeout           int
	Retries           int
	LastStatus        string
	LastStatusChanged time.Time
}

const MonTypeHttp = "http"
const MonStatusActive = "active"
const MonStatusSilenced = "silenced"
const MonLastStatusUp = "up"
const MonLastStatusDown = "down"

func (m Monitor) String() string {
	return fmt.Sprintf("monname:%s type:%s", m.MonName, m.MonType)
}

func GetActiveMonitors() ([]Monitor, error) {
	var err error
	rows, err := Database.Query(fmt.Sprintf("select * from monitor where monstatus=\"%s\" order by monname", MonStatusActive), nil)
	if err != nil {
		log.Printf("failed to query table monitor, error: %s", err)
	}
	defer rows.Close()
	var result []Monitor
	for rows.Next() {
		var id int
		var monname, montype, monstatus, url, expRespCode, laststatus string
		var intrvl, timeout, retries int
		var laststatuschanged time.Time
		err = rows.Scan(&id, &monname, &montype, &monstatus, &url, &intrvl, &expRespCode, &timeout, &retries, &laststatus, &laststatuschanged)
		if err != nil {
			log.Printf("error while scanning monitor table: %s", err)
		}
		result = append(result, Monitor{
			Id:                id,
			MonName:           monname,
			MonType:           montype,
			MonStatus:         monstatus,
			Url:               url,
			Interval:          intrvl,
			ExpRespCode:       expRespCode,
			Timeout:           timeout,
			Retries:           retries,
			LastStatus:        laststatus,
			LastStatusChanged: laststatuschanged,
		})
	}
	return result, err
}

func GetMonitorByName(name string) (Monitor, error) {
	var err error
	var mon Monitor
	selectSQL := "select * from monitor where monname=?"
	statement, err := Database.Prepare(selectSQL)
	if err != nil {
		msg := fmt.Sprintf("failed to prepare stmt for select monitor with name %s, error: %s", name, err)
		log.Print(msg)
		return mon, errors.New(msg)
	} else {
		defer statement.Close()
		var id int
		var monname, montype, monstatus, url, expRespCode, laststatus string
		var intrvl, timeout, retries int
		var laststatuschanged time.Time
		err = statement.QueryRow(name).Scan(&id, &monname, &montype, &monstatus, &url, &intrvl, &expRespCode, &timeout, &retries, &laststatus, &laststatuschanged)
		if err != nil {
			log.Printf("failed to get monitor with name %s, error: %s", monname, err)
			return Monitor{}, err
		}
		return Monitor{
			Id:                id,
			MonName:           monname,
			MonType:           montype,
			MonStatus:         monstatus,
			Url:               url,
			Interval:          intrvl,
			ExpRespCode:       expRespCode,
			Timeout:           timeout,
			Retries:           retries,
			LastStatus:        laststatus,
			LastStatusChanged: laststatuschanged,
		}, err
	}
}

func GetMonitorById(Id int) (Monitor, error) {
	var err error
	var mon Monitor
	selectSQL := "select * from monitor where id=?"
	statement, err := Database.Prepare(selectSQL)
	if err != nil {
		return mon, errors.New(fmt.Sprintf("failed to prepare stmt for select monitor with id %d, error: %s", Id, err))
	} else {
		defer statement.Close()
		var id int
		var monname, montype, monstatus, url, expRespCode, laststatus string
		var intrvl, timeout, retries int
		var laststatuschanged time.Time
		err = statement.QueryRow(Id).Scan(&id, &monname, &montype, &monstatus, &url, &intrvl, &expRespCode, &timeout, &retries, &laststatus, &laststatuschanged)
		if err != nil {
			log.Printf("failed to get monitor with id %d, error: %s", Id, err)
			return Monitor{}, err
		}
		return Monitor{
			Id:                id,
			MonName:           monname,
			MonType:           montype,
			MonStatus:         monstatus,
			Url:               url,
			Interval:          intrvl,
			ExpRespCode:       expRespCode,
			Timeout:           timeout,
			Retries:           retries,
			LastStatus:        laststatus,
			LastStatusChanged: laststatuschanged,
		}, err
	}
}
func UpdateMonitor(mon Monitor) error {
	var err error
	updateSQL := "update monitor set monname=?,montype=?,monstatus=?,url=?,intrvl=?,exp_resp_code=?,timeout=?,retries=?,laststatus=?,laststatuschanged=? where monname=?"
	statement, err := Database.Prepare(updateSQL)
	if err != nil {
		msg := fmt.Sprintf("failed to prepare stmt for update monitor with name %s, error: %s", mon.MonName, err)
		log.Print(msg)
		return errors.New(msg)
	} else {
		defer statement.Close()
		result, err := statement.Exec(mon.MonName, mon.MonType, mon.MonStatus, mon.Url, mon.Interval, mon.ExpRespCode, mon.Timeout, mon.Retries, mon.LastStatus, mon.LastStatusChanged, mon.MonName)
		if err != nil {
			log.Printf("failed to update monitor with name %s, error: %s", mon.MonName, err)
			return err
		}
		numRows, _ := result.RowsAffected()
		if numRows != 1 {
			log.Printf("updated rows for monitor %s is %d (should be 1)", mon.MonName, numRows)
		}
		return err
	}
}
