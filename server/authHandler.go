package main

import (
	"log"
	"net"
	"time"
)

const (
	NoUsername = 0
	Valid      = 1
	WrongPass  = 2
)

func (s *server) handleAuth(conn net.Conn, username, password string) {
	if username == "" || password == "" {
		connErr(conn, "no empty fields allowed for username and password")
		return
	}

	var userStatus int
	err := s.db.QueryRow(`
		SELECT
		CASE
			WHEN NOT EXISTS(SELECT 1 FROM users WHERE username = ?) THEN ?
			WHEN EXISTS(SELECT 1 FROM users WHERE username = ? AND password = ?) THEN ?
			ELSE ?
		END AS userstatus`,
		username, NoUsername, username, password, Valid, WrongPass).
		Scan(&userStatus)
	if err != nil {
		connErr(conn, err.Error())
	}

	s.lock.Lock()
	defer s.lock.Unlock()

	switch userStatus {
	case NoUsername:

		_, err := s.db.Exec(`
			INSERT INTO users (username, password) VALUES (?, ?)`,
			username, password)
		if err != nil {
			connErr(conn, err.Error())
		}
		s.timeMap[username] = time.Now().Format("2006-01-02 15:04:05")
		loginSuccess(conn, "", s.postsList)
		return

	case WrongPass:

		connErr(conn, "invalid password")
		return

	case Valid:

		lastLogin := s.timeMap[username]
		s.timeMap[username] = time.Now().Format("2006-01-02 15:04:05")
		loginSuccess(conn, lastLogin, s.postsList)
		return

	default:
		log.Fatal("This should never be reached")
	}

}
