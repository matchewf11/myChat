package main

import (
	"fmt"
	"net"
)

func (s *server) handleAddRoom(conn net.Conn, username, password, roomName, roomPass string) {

	var validUser bool
	if err := s.db.QueryRow(`
		SELECT EXISTS(
			SELECT 1 
			FROM users 
			WHERE username = ? AND password = ?
	)`, username, password).Scan(&validUser); err != nil {
		connErr(conn, err.Error())
		return
	}

	if !validUser {
		fmt.Println("unvalid user")
		return
	}

	// TODO: add premission to the user fo teh room they created

	if _, err := s.db.Exec(`
		INSERT INTO rooms 
		(room_name, room_password) 
		VALUES (?, ?)
		`, roomName, roomPass); err != nil {
		connErr(conn, err.Error())
		return
	}

	fmt.Println("we just got an add room request")
}

func (s *server) handleDelRoom(conn net.Conn, username, password, roomName string) {

	var validUser bool
	if err := s.db.QueryRow(`
		SELECT EXISTS(
			SELECT 1 
			FROM users 
			WHERE username = ? AND password = ?
	)`, username, password).Scan(&validUser); err != nil {
		connErr(conn, err.Error())
		return
	}

	if !validUser {
		fmt.Println("unvalid user")
		return
	}

	if _, err := s.db.Exec(`
		DELETE FROM rooms
		WHERE room_name = ?
		`, roomName); err != nil {
		connErr(conn, err.Error())
		return
	}

	fmt.Println("we just got an delete room request")
}
