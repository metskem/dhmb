package conf

const DatabaseURL = "file:dhmb.db"
const CreateTablesFile = "resources/sql/create-tables.sql"
const InsertTestDataFile = "resources/sql/insert-testdata.sql"

// Variables to identify the build
var (
	CommitHash string
	VersionTag string
	BuildTime  string
)
