package main

import (
	"encoding/json"
	"fmt"
	"net"
)

// Write an error message across the connection
func connErr(conn net.Conn, body string) {
	sendJSON(conn, map[string]string{
		"status": "error",
		"body":   body,
	})
}

func sendJSON(conn net.Conn, data any) {
	if err := json.NewEncoder(conn).Encode(data); err != nil {
		fmt.Println("error sending to user")
	}
}

func loginSuccess(conn net.Conn, date string, posts []post) {
	sendJSON(conn, struct {
		Status   string `json:"status"`
		Body     string `json:"body"`
		Date     string `json:"date,omitempty"`
		Messages []post `json:"messages"`
	}{
		Status:   "loggedin",
		Body:     "logged in",
		Date:     date,
		Messages: posts,
	})
}
