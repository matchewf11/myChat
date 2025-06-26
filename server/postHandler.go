package main

import (
	"net"
	"time"
)

func (s *server) handlePost(conn net.Conn, username, password, body string) {

	if username == "" || password == "" || body == "" {
		connErr(conn, "username, password, and body must not be empty")
		return
	}

	var valid bool
	if err := s.db.QueryRow(`
		SELECT EXISTS(
			SELECT 1 
			FROM users 
			WHERE username = ? AND password = ?
		)`, username, password).Scan(&valid); err != nil {
		connErr(conn, err.Error())
		return
	}

	if !valid {
		connErr(conn, "invalid username or password")
		return
	}

	if _, err := s.db.Exec(`
		INSERT INTO posts (body, author)
		VALUES (?, (SELECT id FROM users WHERE username = ?))
		`, body, username); err != nil {
		connErr(conn, err.Error())
		return
	}

	s.lock.Lock()

	conns := make([]net.Conn, 0, len(s.usersMap))
	for user := range s.usersMap {
		if user != conn {
			conns = append(conns, user)
		}
	}
	s.lock.Unlock()

	for _, user := range conns {
		if err := sendJSON(user, map[string]string{
			"status":   "received",
			"username": username,
			"body":     body,
			"date":     time.Now().Format("2006-01-02 15:04:05"),
		}); err != nil {
			connErr(conn, err.Error())
			return
		}
	}
	sendJSON(conn, map[string]string{
		"status": "sent",
		"body":   "you sent a post request",
	})
}
