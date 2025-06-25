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
			INSERT INTO users (username, password, last_login) VALUES (?, ?, ?)`,
			username, password, time.Now().Format("2006-01-02 15:04:05"))
		if err != nil {
			connErr(conn, err.Error())
		}
		loginSuccess(conn, "", s.postsList)
		return

	case WrongPass:

		connErr(conn, "invalid password")
		return

	case Valid:

		var lastLogin string

		err := s.db.QueryRow(`
			SELECT last_login
			FROM users
			WHERE username = ?
			`, username).Scan(&lastLogin)
		if err != nil {
			connErr(conn, err.Error())
			return
		}

		loginSuccess(conn, lastLogin, s.postsList)

		_, err = s.db.Exec(`
			UPDATE users
			SET last_login = ?
			WHERE username = ?
			`, time.Now().Format("2006-01-02 15:04:05"))
		if err != nil {
			connErr(conn, err.Error())
			return
		}

		return

	default:
		log.Fatal("This should never be reached")
	}

}
