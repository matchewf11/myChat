package main

import (
	"fmt"
	"net"
	"time"
)

func (s *server) handleAuth(conn net.Conn, username, password string) {
	if username == "" || password == "" {
		connErr(conn, "no empty fields allowed for username and password")
		return
	}

	s.lock.Lock()

	fmt.Println("ALl the usernames:")
	for key := range s.passwordMap {
		fmt.Println(key)
	}

	val, exists := s.passwordMap[username]

	if !exists {
		s.passwordMap[username] = password
		s.timeMap[username] = time.Now().Format("2006-01-02 15:04:05")
		s.lock.Unlock()
		loginSuccess(conn, "", s.postsList)
		fmt.Println("username does not exist")
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

	fmt.Println(lastLogin)

	loginSuccess(conn, lastLogin, s.postsList)
}
