package check

import (
	"github.com/metskem/dhmb/db"
	"log"
)

/*
  Iterate over the monitors in the monitor table, and start a separate thread for a check
*/

func Runner() {
	for ix, m := range db.GetMonitors() {
		index := ix
		monitor := m
		go func(db.Monitor) {
			log.Printf("background monitor %d: %s", index, monitor)
			if monitor.MonType == db.MonTypeHttp {
				Loop(monitor)
			}
		}(monitor)
	}
}
