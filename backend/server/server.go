package server

import (
	"bufio"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"sync"
)

type userStatus int

const (
	NoUsername userStatus = iota
	WrongPass
	Valid
	Invalid
)

const (
	loginMethod = "LOGIN"
	postMethod  = "POST"
	roomMethod  = "ROOM"
)

type server struct {
	incoming  chan message
	db        *sql.DB
	clients   map[net.Conn]bool
	clientsMu sync.Mutex
}

type message struct {
	Method   string `json:"method"`
	Username string `json:"username"`
	Password string `json:"password"`
	RoomName string `json:"room_name"`
	RoomPass string `json:"room_pass"`
	Body     string `json:"body"`
	conn     net.Conn
}

type post struct {
	Author string `json:"author"`
	Body   string `json:"body"`
	Date   string `json:"date,omitempty"`
}

type room struct {
	Name string `json:"room_name"`
}

func InitServer(db *sql.DB) *server {
	return &server{
		db:       db,
		incoming: make(chan message),
		clients:  make(map[net.Conn]bool),
	}
}

func (s *server) ListenConnection(conn net.Conn) {

	s.clientsMu.Lock()
	s.clients[conn] = true
	s.clientsMu.Unlock()

	defer func() {
		s.clientsMu.Lock()
		delete(s.clients, conn)
		s.clientsMu.Unlock()
		conn.Close()
	}()

	scanner := bufio.NewScanner(conn)

	for scanner.Scan() {
		var req message
		if err := json.Unmarshal([]byte(scanner.Text()), &req); err != nil {
			sendErr(conn, "unable to parse into json")
			continue
		}
		req.conn = conn
		s.incoming <- req
	}
}

func (s *server) HandleChan() { // Finish this!!!!!!!

	for {
		msg := <-s.incoming

		currStatus := s.validateUser(msg.Username, msg.Password)

		switch currStatus {
		case Invalid, WrongPass:
			sendErr(msg.conn, "invalid username or password")
		case NoUsername:
			if msg.Method == loginMethod {
				s.makeNewUser(msg.Username, msg.Password)
				s.sendLoginInfo(msg.conn, "signed up")
			} else {
				sendErr(msg.conn, "still need to sign in")
			}
		case Valid:
			switch msg.Method {
			case loginMethod:
				s.sendLoginInfo(msg.conn, "logged in")
			case postMethod:
				s.makePost(msg.Username, msg.Body)
				s.broadcast(map[string]any{
					"status": "new_post",
					"post": post{
						Author: msg.Username,
						Body:   msg.Body,
					},
				})
			case roomMethod:
				s.makeRoom(msg.conn, msg.RoomName)
				s.broadcast(map[string]any{
					"status": "new_room",
					"room":   room{Name: msg.RoomName},
				})
			default:
				sendErr(msg.conn, "not valid method")
			}
		}
	}
}

func (s *server) broadcast(msg any) {
	s.clientsMu.Lock()
	defer s.clientsMu.Unlock()

	for conn := range s.clients {
		if err := json.NewEncoder(conn).Encode(msg); err != nil {
			log.Printf("error sending to %v: %v\n", conn.RemoteAddr(), err)
		}
	}
}

func (s *server) makePost(username, body string) {
	_, err := s.db.Exec(`
		INSERT INTO posts (author, body)
		VALUES ((SELECT id FROM users WHERE username = ?), ?)
	`, username, body)
	crashIfErr(err)
}

func (s *server) makeRoom(conn net.Conn, roomName string) {
	if len(roomName) > 16 {
		sendErr(conn, "invalid room name")
		return
	}
	_, err := s.db.Exec(`
		INSERT INTO rooms (room_name)
		VALUES (?)
	`, roomName)
	crashIfErr(err)
}

func (s *server) sendLoginInfo(conn net.Conn, status string) {
	sendMsg(conn, map[string]any{
		"status": status,
		"posts":  s.getPosts(),
		"rooms":  s.getRooms(),
	})
}

func (s *server) getRooms() []room {
	roomQry, err := s.db.Query(`SELECT name FROM rooms`)
	crashIfErr(err)
	defer roomQry.Close()

	roomList := make([]room, 0)
	for roomQry.Next() {
		var name string
		err = roomQry.Scan(&name)
		crashIfErr(err)
		roomList = append(roomList, room{Name: name})
	}
	return roomList
}

func (s *server) getPosts() []post {
	postQry, err := s.db.Query(`
			SELECT posts.body, users.username, posts.created_at
			FROM posts
			JOIN users ON posts.author = users.id
			ORDER BY posts.created_at
	`)
	crashIfErr(err)
	defer postQry.Close()

	postList := make([]post, 0)
	for postQry.Next() {
		var body, author, created_at string
		err = postQry.Scan(&body, &author, &created_at)
		crashIfErr(err)
		postList = append(postList, post{
			Body:   body,
			Author: author,
			Date:   created_at,
		})

	}

	return postList
}

func (s *server) validateUser(username, password string) userStatus {

	if username == "" || password == "" || len(username) > 16 {
		return Invalid
	}

	var currStatus userStatus
	if err := s.db.QueryRow(`
		SELECT CASE
			WHEN NOT EXISTS(SELECT 1 FROM users WHERE username = ?) THEN ?
			WHEN EXISTS(SELECT 1 FROM users WHERE username = ? AND password = ?) THEN ?
			ELSE ?  
		END AS userstatus`,
		username, NoUsername, username, password, Valid, WrongPass).
		Scan(&currStatus); err != nil {
		return Invalid
	}

	return currStatus
}

func sendErr(conn net.Conn, errMsg string) {
	if err := json.NewEncoder(conn).Encode(map[string]string{
		"error": errMsg,
	}); err != nil {
		fmt.Println("can send err message to user")
	}
}

func sendMsg(conn net.Conn, msg any) {
	err := json.NewEncoder(conn).Encode(msg)
	crashIfErr(err)
}

func (s *server) makeNewUser(username, password string) {
	_, err := s.db.Exec(`
		INSERT INTO users (username, password)
		VALUES (?, ?)
	`, username, password)
	crashIfErr(err)
}

func crashIfErr(err error) {
	if err != nil {
		log.Fatal(err.Error())
	}
}
