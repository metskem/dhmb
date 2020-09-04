package db

import (
	"fmt"
	"log"
)

type Monitor struct {
	id          int
	monname     string
	montype     string
	url         string
	interval    int
	expRespCode int
	timeout     int
}

func (m Monitor) String() string {
	return fmt.Sprintf("monname:%s type:%s", m.monname, m.montype)
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
		var url string
		var intrvl int
		var expRespCode int
		var timeout int
		err = rows.Scan(&id, &monname, &montype, &url, &intrvl, &expRespCode, &timeout)
		if err != nil {
			log.Printf("error while scanning monitor table: %s", err)
		}
		result = append(result, Monitor{
			id:          id,
			monname:     monname,
			montype:     montype,
			url:         url,
			interval:    intrvl,
			expRespCode: expRespCode,
			timeout:     timeout,
		})
	}
	return result
}
