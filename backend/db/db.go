package db

import (
	_ "github.com/tursodatabase/go-libsql"

	"database/sql"

	_ "embed"
)

//go:embed createTables.sql
var createTableSql string

func InitDB() (*sql.DB, error) {
	db, err := sql.Open("libsql", "db/myChatDb")
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
