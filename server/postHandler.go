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
	err := s.db.QueryRow(`
		SELECT EXISTS(
			SELECT 1 
			FROM users 
			WHERE username = ? AND password = ?
		)`, username, password).Scan(&valid)

	if err != nil {
		connErr(conn, err.Error())
		return
	}

	if !valid {
		connErr(conn, "invalid username or password")
		return
	}

	currTime := time.Now().Format("2006-01-02 15:04:05")
	newPost := post{
		Body:   body,
		Date:   currTime,
		Author: username,
	}

	s.lock.Lock()
	s.postsList = append(s.postsList, newPost)

	conns := make([]net.Conn, 0, len(s.usersMap))
	for user := range s.usersMap {
		if user != conn {
			conns = append(conns, user)
		}
	}
	s.lock.Unlock()

	for _, user := range conns {
		sendJSON(user, map[string]string{
			"status":   "received",
			"username": username,
			"body":     body,
			"date":     currTime,
		})
	}

	sendJSON(conn, map[string]string{
		"status": "sent",
		"body":   "you sent a post request",
	})

}
