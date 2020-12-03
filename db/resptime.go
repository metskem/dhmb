package db

import (
	"fmt"
	"github.com/metskem/dhmb/conf"
	"log"
	"time"
)

type RespTime struct {
	Id        int
	Timestamp time.Time
	MonId     int
	Time      int64
}

func (respTime RespTime) String() string {
	return fmt.Sprintf("Id:%d monid:%d resptime:%d", respTime.Id, respTime.MonId, respTime.Time)
}

func GetRespTimes() []RespTime {
	rows, err := Database.Query("select id, timestamp, monid, time from resptime", nil)
	if err != nil {
		log.Printf("failed to query table resptime, error: %s", err)
	}
	defer rows.Close()
	var result []RespTime
	for rows.Next() {
		var id int
		var timestamp time.Time
		var monId int
		var time int64
		err = rows.Scan(&id, &timestamp, &monId, &time)
		if err != nil {
			log.Printf("error while scanning resptime table: %s", err)
		}
		result = append(result, RespTime{
			Id:        id,
			Timestamp: timestamp,
			MonId:     monId,
			Time:      time,
		})
	}
	return result
}

/**
  Insert a row into the resptime table using the given resptime. Returns the lastInsertId of the insert operation.
*/
func InsertRespTime(respTime RespTime) int64 {
	insertSQL := "insert into resptime(timestamp, monid, time) values(?,?,?)"
	statement, err := Database.Prepare(insertSQL)
	defer statement.Close()
	if err != nil {
		log.Printf("failed to prepare stmt for insert into resptime, error: %s", err)
		return 0
	} else {
		result, err := statement.Exec(respTime.Timestamp, respTime.MonId, respTime.Time)
		if err != nil {
			log.Printf("failed to insert monid %d, timestamp %s, time %d, error: %s", respTime.Timestamp.Format(time.RFC3339), respTime.MonId, respTime.Time, err)
			return 0
		} else {
			lastInsertId, err := result.LastInsertId()
			if err == nil {
				return lastInsertId
			} else {
				log.Printf("no resptime row was inserted, err: %s", err)
				return 0
			}
		}
	}
}

func CleanupOldStuffForMonitor(m Monitor) int64 {
	deleteSQL := `delete from resptime where monid=(select id from monitor where monname=?) and id not in (select r.id from resptime r,monitor m where m.id=r.monid and m.monname=? order by timestamp desc limit ?)`
	statement, err := Database.Prepare(deleteSQL)
	defer statement.Close()
	if err != nil {
		log.Printf("failed to prepare stmt for delete from resptime, error: %s", err)
		return 0
	} else {
		result, err := statement.Exec(m.MonName, m.MonName, conf.MaxRowsResptime)
		if err != nil {
			log.Printf("failed to delete rows from resptime for monitor %s, error: %s", m.MonName, err)
			return 0
		} else {
			numDeleted, err := result.RowsAffected()
			if err == nil {
				log.Printf("%d rows deleted from resptime table for monitor %s", numDeleted, m.MonName)
				return numDeleted
			} else {
				log.Printf("no resptime rows were deleted, err: %s", err)
				return 0
			}
		}
	}
}
