package main

import (
	"log"
	"net"
	"time"
)

type userStatus int

const (
	NoUsername userStatus = iota
	Valid
	WrongPass
)

func (s *server) handleAuth(conn net.Conn, username, password string) {

	if username == "" || password == "" {
		connErr(conn, "no empty fields allowed for username and password")
		return
	}

	var userStatus userStatus
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

	switch userStatus {
	case NoUsername:
		// TODO: Fix this later
		_, err := s.db.Exec(`
			INSERT INTO users (username, password, last_login) VALUES (?, ?, ?)`,
			username, password, time.Now().Format("2006-01-02 15:04:05"))
		if err != nil {
			connErr(conn, err.Error())
			return
		}
		postList, err := s.getPosts()
		if err != nil {
			connErr(conn, err.Error())
			return
		}
		if err = loginSuccess(conn, "", postList); err != nil {
			connErr(conn, err.Error())
			return
		}
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
		postList, err := s.getPosts()
		if err != nil {
			connErr(conn, err.Error())
			return
		}
		if err = loginSuccess(conn, lastLogin, postList); err != nil {
			connErr(conn, err.Error())
			return
		}
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
