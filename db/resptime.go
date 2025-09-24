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

func GetLatestRespTimesByMonname(monname string) []RespTime {
	var result []RespTime
	rows, err := Database.Query("select * from (select r.id, r.timestamp, r.monid, r.time from resptime r, monitor m where r.monid=m.id and m.monname=? order by r.timestamp desc limit ?) order by timestamp", monname, conf.MaxPlots)
	if err != nil {
		log.Printf("failed to query table resptime, error: %s", err)
	} else {
		defer rows.Close()
		for rows.Next() {
			var id, monId int
			var timestamp time.Time
			var respTime int64
			err = rows.Scan(&id, &timestamp, &monId, &respTime)
			if err != nil {
				log.Printf("error while scanning resptime table: %s", err)
			}
			result = append(result, RespTime{
				Id:        id,
				Timestamp: timestamp,
				MonId:     monId,
				Time:      respTime,
			})
		}
	}
	return result
}

// GetNewestTimestamps - return the last resptime update for each monitor, a map of timestamps is returned, key-ed by monid
func GetNewestTimestamps() map[int]time.Time {
	var dateFormat = "2006-01-02 15:04:05Z07:00"
	var result = make(map[int]time.Time)
	rows, err := Database.Query(fmt.Sprintf(`select monid,max(timestamp),monstatus from resptime where monstatus!="%s" group by monid order by monid`, MonStatusInactive))
	if err != nil {
		log.Printf("failed to get newest timestamps from resptime, error: %s", err)
	} else {
		defer rows.Close()
		for rows.Next() {
			var monId int
			var timestampStr string
			err = rows.Scan(&monId, &timestampStr)
			if err != nil {
				log.Printf("error while scanning resptime table: %s", err)
			}
			timestamp, err := time.Parse(dateFormat, timestampStr)
			if err != nil {
				log.Printf("failed to parse timestamp from resptime, error: %s", err)
			}
			result[monId] = timestamp
		}
	}
	return result
}

// InsertRespTime - Insert a row into the resptime table using the given resptime. Returns the lastInsertId of the insert operation. */
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
			log.Printf("failed to insert monid %d, timestamp %s, time %d, error: %s", respTime.MonId, respTime.Timestamp.Format(time.RFC3339), respTime.Time, err)
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
