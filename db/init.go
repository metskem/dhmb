package db

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"github.com/metskem/dhmb/conf"
	"io/ioutil"
	"log"
	"os"
)

func Initdb() *sql.DB {
	var DbExists bool
	if _, err := os.Stat(conf.DatabaseURL); err == nil {
		log.Printf("database already exists, opening it...\n")
		DbExists = true
	} else {
		log.Printf("database does not yet exist, creating it...\n")
		DbExists = false
	}

	db, err := sql.Open("sqlite3", conf.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}

	if !DbExists {
		sqlStmts, err := ioutil.ReadFile(conf.CreateTablesFile)
		if err != nil {
			log.Fatal(err)
		}
		_, err = db.Exec(string(sqlStmts))
		if err != nil {
			log.Fatalf("%q: %s\n", err, sqlStmts)
		}

		sqlStmts, err = ioutil.ReadFile(conf.InsertTestDataFile)
		if err != nil {
			log.Fatal(err)
		}
		_, err = db.Exec(string(sqlStmts))
		if err != nil {
			log.Fatalf("%q: %s\n", err, sqlStmts)
		}
	}

	// show the # of rows in the tables
	rows, err := db.Query("select * from monitor", nil)
	if err != nil {
		db.Close()
		log.Fatalf("failed to query table %, error: %s", "monitor", err)
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		var monname string
		var montype string
		var url string
		var intrvl int
		err = rows.Scan(&id, &monname, &montype, &url, &intrvl)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("%d,%s,%s,%s,%d\n", id, monname, montype, url, intrvl)
	}
	return db
}
