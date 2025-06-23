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
	lock     sync.Mutex
	users    map[net.Conn]bool
	userAuth map[string]string
	userTime map[string]string
	messages []message
}

type message struct {
	Body   string `json:"body"`
	Author string `json:"author"`
	Date   string `json:"date"`
}

func main() {

	serv := &server{
		users:    make(map[net.Conn]bool),
		userAuth: make(map[string]string),
		userTime: make(map[string]string),
		messages: make([]message, 0),
	}

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
			Password string `json:"password"`
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
		case "AUTH":

			fmt.Println(req.Username, req.Password)

			if req.Username == "" || req.Password == "" {
				if err := json.NewEncoder(conn).Encode(map[string]string{
					"status": "login fail",
					"body":   "does not allow for empty fields",
				}); err != nil {
					fmt.Println("error sending to user")
				}
				break
			}

			s.lock.Lock()
			val, contains := s.userAuth[req.Username]
			s.lock.Unlock()

			if !contains {
				s.lock.Lock()
				s.userAuth[req.Username] = req.Password
				s.lock.Unlock()

				// TODO: ALso send old messages here

				if err := json.NewEncoder(conn).Encode(map[string]string{
					"status": "loggedin",
					"body":   "logged in",
				}); err != nil {
					fmt.Println("error sending to user")
				}
				break
			}

			if val != req.Password {
				if err := json.NewEncoder(conn).Encode(map[string]string{
					"status": "login fail",
					"body":   "invalid password",
				}); err != nil {
					fmt.Println("error sending to user")
				}
				break
			}

			var date string = ""

			s.lock.Lock()

			lastLogin, contains := s.userTime[req.Username]
			if contains {
				date = lastLogin
			}
			s.userTime[req.Username] = time.Now().Format("2006-01-02 15:04:05")
			s.lock.Unlock()

			jsonMessages, err := json.Marshal(s.messages)
			if err != nil {
				log.Fatal("can't marshal jsonMessages")
			}

			fmt.Println(s.messages)

			// TODO: properly format the messages part of the json

			if err := json.NewEncoder(conn).Encode(map[string]string{
				"status":   "loggedin",
				"body":     "logged in",
				"date":     date,
				"messages": string(jsonMessages),
			}); err != nil {
				fmt.Println("error sending to user")
			}

		case "POST":

			s.lock.Lock()
			val, has := s.userAuth[req.Username]
			s.lock.Unlock()

			if !has {
				if err := json.NewEncoder(conn).Encode(map[string]string{
					"status": "error",
					"body":   "invalid username",
				}); err != nil {
					fmt.Println("could not send sucess message")
				}
				break
			}

			if val != req.Password {
				if err := json.NewEncoder(conn).Encode(map[string]string{
					"status": "error",
					"body":   "invalid password",
				}); err != nil {
					fmt.Println("could not send sucess message")
				}
				break
			}

			s.lock.Lock()

			s.messages = append(s.messages, message{
				Body:   req.Body,
				Date:   time.Now().Format("2006-01-02 15:04:05"),
				Author: req.Username,
			})

			for user := range s.users {
				if user == conn {
					continue
				}
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
				"status": "sent",
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
