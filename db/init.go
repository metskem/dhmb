package db

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"github.com/metskem/dhmb/conf"
	"io/ioutil"
	"log"
	"os"
)

var Database *sql.DB

func Initdb() {
	var DbExists bool
	var err error
	if _, err = os.Stat(conf.DatabaseURL); err == nil {
		log.Printf("database already exists, opening it...\n")
		DbExists = true
	} else {
		log.Printf("database does not yet exist, creating it...\n")
		DbExists = false
	}

	Database, err = sql.Open("sqlite3", conf.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}

	if !DbExists {
		sqlStmts, err := ioutil.ReadFile(conf.CreateTablesFile)
		if err != nil {
			log.Fatal(err)
		}
		_, err = Database.Exec(string(sqlStmts))
		if err != nil {
			log.Fatalf("%q: %s\n", err, sqlStmts)
		}

		sqlStmts, err = ioutil.ReadFile(conf.InsertTestDataFile)
		if err != nil {
			log.Fatal(err)
		}
		_, err = Database.Exec(string(sqlStmts))
		if err != nil {
			log.Fatalf("%q: %s\n", err, sqlStmts)
		}
	}
}
