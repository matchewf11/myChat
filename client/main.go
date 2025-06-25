package main

import (
	"log"

	"github.com/rivo/tview"
)

const (
	StatusLoginFail = "login fail"
	StatusLoggedIn  = "loggedin"
)

func main() {

	view := initView()
	defer view.conn.Close()

	loginForm := view.getLoginForm()
	roomList := view.getRoomList()
	textFlex, inputArea, textView := view.genTextArea()

	rowFlex := tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(loginForm, 0, 1, true).
		AddItem(roomList, 0, 0, false).
		AddItem(textFlex, 0, 0, false)

	go listenServer(view, loginForm, textView, inputArea, rowFlex, roomList, textFlex)

	if err := view.app.SetRoot(rowFlex, true).EnableMouse(true).Run(); err != nil {
		log.Fatal(err)
	}
}
