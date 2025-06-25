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

func (s *server) getPosts() ([]post, error) {

	postList := make([]post, 0)

	postQry, err := s.db.Query(`
			SELECT posts.body, users.username, posts.created_at
			FROM posts
			JOIN users ON posts.author = users.id
			ORDER BY posts.created_at
		`)
	if err != nil {
		return nil, err
	}
	defer postQry.Close()

	for postQry.Next() {
		var body, author, created_at string
		err = postQry.Scan(&body, &author, &created_at)
		if err != nil {
			return nil, err
		}
		postList = append(postList, post{
			Body:   body,
			Author: author,
			Date:   created_at,
		})
	}

	return postList, nil

}
