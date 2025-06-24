package main

import (
	"net"
	"time"
)

func (s *server) handleAuth(conn net.Conn, username, password string) {
	if username == "" || password == "" {
		connErr(conn, "no empty fields allowed for username and password")
		return
	}

	s.lock.Lock()
	val, exists := s.passwordMap[username]

	if !exists {
		s.passwordMap[username] = password
		s.lock.Unlock()
		loginSuccess(conn, "", s.postsList)
		return
	}

	if val != password {
		s.lock.Unlock()
		connErr(conn, "invalid password")
		return
	}

	lastLogin := s.timeMap[username]
	s.timeMap[username] = time.Now().Format("2006-01-02 15:04:05")
	s.lock.Unlock()

	loginSuccess(conn, lastLogin, s.postsList)
}
