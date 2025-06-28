package db

import (
	"database/sql"

	_ "embed"
	_ "github.com/mattn/go-sqlite3"
)

//go:embed createTables.sql
var createTableSql string

func InitDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "db/myChatDb")
	if err != nil {
		return nil, err
	}

	if _, err = db.Exec("PRAGMA foreign_keys = ON;"); err != nil {
		return nil, err
	}

	if _, err = db.Exec(createTableSql); err != nil {
		return nil, err
	}

	return db, nil
}
