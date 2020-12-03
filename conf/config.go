package conf

import (
	"log"
	"os"
	"strconv"
)

const DatabaseURL = "file:dhmb.db"
const CreateTablesFile = "resources/sql/create-tables.sql"
const InsertTestDataFile = "resources/sql/insert-testdata.sql"

// Variables to identify the build
var (
	CommitHash string
	VersionTag string
	BuildTime  string

	BotToken           = os.Getenv("BOT_TOKEN")
	DebugStr           = os.Getenv("DEBUG")
	Debug              bool
	MaxRowsResptimeStr = os.Getenv("MAX_ROWS_RESPTIME")
	MaxRowsResptime    int
)

func EnvironmentComplete() bool {
	envComplete := true

	if len(BotToken) == 0 {
		log.Print("missing envvar \"BOT_TOKEN\"")
		envComplete = false
	}

	Debug = false
	if DebugStr == "true" {
		Debug = true
	}

	MaxRowsResptime = 1000
	if MaxRowsResptimeStr != "" {
		MaxRowsResptime, _ = strconv.Atoi(MaxRowsResptimeStr)
	}

	return envComplete
}
