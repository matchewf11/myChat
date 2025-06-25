package main

import (
	"fmt"
	"net"
)

func main() {

	ln, err := net.Listen("tcp", ":9000")
	if err != nil {
		fmt.Println("err creating server:", err)
	}
	defer ln.Close()

	db := initDb()
	defer db.Close()

	serv := newServer(db)

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
