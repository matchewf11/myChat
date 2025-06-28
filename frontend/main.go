package main

import (
	"log"
	"myChat/frontend/view"

	"github.com/rivo/tview"
)

func main() {
	v := view.InitView()
	defer v.Conn.Close()

	loginForm := v.GetLoginForm()
	roomFlex, roomList := v.GetRoomList()
	textFlex, inputArea, textView := v.GenTextArea()

	rowFlex := tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(loginForm, 0, 1, true).
		AddItem(roomFlex, 0, 0, false).
		AddItem(textFlex, 0, 0, false)

	go view.ListenServer(v, loginForm, textView, inputArea, rowFlex, roomList, textFlex, roomFlex)

	if err := v.App.SetRoot(rowFlex, true).EnableMouse(true).Run(); err != nil {
		log.Fatal(err)
	}

}
