package main

import (
	_ "embed"
	"fmt"
	"log"
	"myChat/backend/db"
	"myChat/backend/server"
	"net"
)

// Starting point of the server
func main() {

	ln, err := net.Listen("tcp", ":9000")
	if err != nil {
		log.Fatal(err.Error())
	}
	defer ln.Close()

	db, err := db.InitDB()
	if err != nil {
		log.Fatal(err.Error())
	}
	defer db.Close()

	svr := server.InitServer(db)
	go svr.HandleChan()

	fmt.Println("Server listening @ :9000")

	for {
		if conn, err := ln.Accept(); err != nil {
			log.Fatal(err)
		} else {
			go svr.ListenConnection(conn)
		}
	}
}
