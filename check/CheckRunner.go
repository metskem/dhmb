package check

import (
	"github.com/metskem/dhmb/db"
	"log"
)

/*
  Iterate over the monitors in the monitor table, and start a separate thread for a check
*/

func CheckRunner() {
	for ix, monitor := range db.GetMonitors() {
		log.Printf("found monitor %d : %s", ix, monitor)
	}
}
