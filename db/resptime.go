package db

import (
	"fmt"
	"log"
)

type RespTime struct {
	Id    int
	MonId int
	Time  int64
}

func (respTime RespTime) String() string {
	return fmt.Sprintf("Id:%d monid:%d resptime:%d", respTime.Id, respTime.MonId, respTime.Time)
}

func GetRespTimes() []RespTime {
	rows, err := Database.Query("select * from resptime", nil)
	if err != nil {
		log.Printf("failed to query table resptime, error: %s", err)
	}
	defer rows.Close()
	var result []RespTime
	for rows.Next() {
		var id int
		var monId int
		var time int64
		err = rows.Scan(&id, &monId, &time)
		if err != nil {
			log.Printf("error while scanning resptime table: %s", err)
		}
		result = append(result, RespTime{
			Id:    id,
			MonId: monId,
			Time:  time,
		})
	}
	return result
}

/**
  Insert a row into the resptime table using the given resptime. Returns the lastInsertId of the insert operation.
*/
func InsertRespTime(respTime RespTime) int64 {
	insertSQL := "insert into resptime(monid, time) values(?,?)"
	statement, err := Database.Prepare(insertSQL)
	defer statement.Close()
	if err != nil {
		log.Printf("failed to prepare stmt for insert into resptime, error: %s", err)
		return 0
	} else {
		result, err := statement.Exec(respTime.MonId, respTime.Time)
		if err != nil {
			log.Printf("failed to insert monid %d, time %d, error: %s", respTime.MonId, respTime.Time, err)
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
