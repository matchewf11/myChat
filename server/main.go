package main

import (
	"fmt"
	"net"
)

// Starting point of the server
func main() {

	ln, err := net.Listen("tcp", ":9000")
	if err != nil {
		fmt.Println("err creating server:", err)
	}
	defer ln.Close()

	fmt.Println("Server listening @ :9000")

	serv := newServer()

	for {
		if conn, err := ln.Accept(); err != nil {
			fmt.Println("accept error:", err)
			continue
		} else {
			go serv.handleConn(conn)
		}
	}
}
