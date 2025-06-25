package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
)

func initDb() *sql.DB {

	db, err := sql.Open("sqlite3", "db/myChatDb")
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec("PRAGMA foreign_keys = ON;")
	if err != nil {
		log.Fatal(err)
	}

	contentBytes, err := os.ReadFile("db/createTables.sql")
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(string(contentBytes))
	if err != nil {
		log.Fatal(err)
	}

	return db
}
