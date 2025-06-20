package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"sync"
	"time"
)

type server struct {
	users map[net.Conn]bool
	lock  sync.Mutex
}

func main() {

	serv := &server{users: make(map[net.Conn]bool)}

	ln, err := net.Listen("tcp", ":9000")
	if err != nil {
		log.Fatal(err)
	}
	defer ln.Close()

	fmt.Println("Server listening @ :9000")

	for {
		if conn, err := ln.Accept(); err != nil {
			fmt.Println("accept error:", err)
			continue
		} else {
			go serv.handleConn(conn)
		}
	}
}

func (s *server) handleConn(conn net.Conn) {

	s.lock.Lock()
	s.users[conn] = true
	s.lock.Unlock()
	fmt.Printf("opened connection %v\n", conn)

	defer func() {
		if r := recover(); r != nil {
			log.Printf("recovered crashed go routine")
		}
		conn.Close()
		s.lock.Lock()
		delete(s.users, conn)
		s.lock.Unlock()
		fmt.Printf("closed connection %v\n", conn)
	}()

	scanner := bufio.NewScanner(conn)

	for scanner.Scan() {

		var req struct {
			Method   string `json:"method"`
			Username string `json:"username"`
			Body     string `json:"body"`
		}

		if err := json.Unmarshal([]byte(scanner.Text()), &req); err != nil {
			json.NewEncoder(conn).Encode(map[string]string{
				"status": "error",
				"body":   "invalid json",
			})
			continue
		}

		switch req.Method {
		case "POST":

			s.lock.Lock()

			for user := range s.users {

				if err := json.NewEncoder(user).Encode(map[string]string{
					"status":   "recieved",
					"username": req.Username,
					"body":     req.Body,
					"date":     time.Now().Format("2006-01-02 15:04:05"),
				}); err != nil {
					fmt.Println("error sending to user")
				}

			}

			s.lock.Unlock()

			fmt.Println(req.Body)

			if err := json.NewEncoder(conn).Encode(map[string]string{
				"status": "success",
				"body":   "you sent a post request",
			}); err != nil {
				fmt.Println("could not send sucess message")
			}

		default:
			if err := json.NewEncoder(conn).Encode(map[string]string{
				"status": "error",
				"body":   "does not support method",
			}); err != nil {
				fmt.Println("could not send error message")
			}
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Connection error:", err)
		return
	}

}
