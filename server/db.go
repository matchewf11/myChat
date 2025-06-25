package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
)

func initDb() *sql.DB {

	db, err := sql.Open("sqlite3", "./db/myChatDb")
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec("PRAGMA foreign_keys = ON;")
	if err != nil {
		log.Fatal(err)
	}

	filePath := "./db/createTables.sql"

	contentBytes, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}

	fileContentString := string(contentBytes)

	_, err = db.Exec(fileContentString)
	if err != nil {
		log.Fatal(err)
	}

	return db
}
