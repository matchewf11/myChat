package main

import (
	"bufio"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"sync"
)

type server struct {
	lock     sync.Mutex
	usersMap map[net.Conn]bool
	db       *sql.DB
}

type post struct {
	Body   string `json:"body"`
	Author string `json:"author"`
	Date   string `json:"date"`
}

func newServer(db *sql.DB) *server {
	return &server{
		usersMap: make(map[net.Conn]bool),
		db:       db,
	}
}

func (s *server) handleConn(conn net.Conn) {

	s.lock.Lock()
	s.usersMap[conn] = true
	s.lock.Unlock()

	fmt.Printf("opened connection %v\n", conn)

	defer func() {
		// TODO: change last login to curr time
		if r := recover(); r != nil {
			log.Printf("recovered crashed go routine")
		}
		conn.Close()
		s.lock.Lock()
		delete(s.usersMap, conn)
		s.lock.Unlock()
		fmt.Printf("closed connection %v\n", conn)
	}()

	scanner := bufio.NewScanner(conn)

	for scanner.Scan() {

		var req struct {
			Method   string `json:"method"`
			Username string `json:"username"`
			Password string `json:"password"`
			RoomName string `json:"room_name"`
			RoomPass string `json:"room_pass"`
			Body     string `json:"body"`
		}

		if err := json.Unmarshal([]byte(scanner.Text()), &req); err != nil {
			connErr(conn, "unable to parse json")
			continue
		}

		switch req.Method {
		case "AUTH":
			s.handleAuth(conn, req.Username, req.Password)
		case "POST":
			s.handlePost(conn, req.Username, req.Password, req.Body)
		case "ADD_ROOM":
			s.handleAddRoom(conn, req.Username, req.Password, req.RoomName, req.RoomPass)
		case "DEL_ROOM":
			s.handleDelRoom(conn, req.Username, req.Password, req.RoomName)
		default:
			connErr(conn, "invalid body")
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Connection error:", err)
		return
	}

}
