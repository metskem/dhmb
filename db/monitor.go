package db

import (
	"fmt"
	"log"
)

type Monitor struct {
	Id          int
	MonName     string
	MonType     string
	MonStatus   string
	Url         string
	Interval    int
	ExpRespCode int
	Timeout     int
	Retries     int
}

const MonTypeHttp = "http"
const MonStatusActive = "active"
const MonStatusInactive = "inactive"

func (m Monitor) String() string {
	return fmt.Sprintf("monname:%s type:%s", m.MonName, m.MonType)
}

func GetMonitors() []Monitor {
	rows, err := Database.Query("select * from monitor order by monname", nil)
	if err != nil {
		log.Fatalf("failed to query table monitor, error: %s", err)
	}
	defer rows.Close()
	var result []Monitor
	for rows.Next() {
		var id int
		var monname string
		var montype string
		var monstatus string
		var url string
		var intrvl int
		var expRespCode int
		var timeout int
		var retries int
		err = rows.Scan(&id, &monname, &montype, &monstatus, &url, &intrvl, &expRespCode, &timeout, &retries)
		if err != nil {
			log.Printf("error while scanning monitor table: %s", err)
		}
		result = append(result, Monitor{
			Id:          id,
			MonName:     monname,
			MonType:     montype,
			MonStatus:   monstatus,
			Url:         url,
			Interval:    intrvl,
			ExpRespCode: expRespCode,
			Timeout:     timeout,
			Retries:     retries,
		})
	}
	return result
}
