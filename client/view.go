package main

import (
	"fmt"
	"log"
	"net"

	"github.com/rivo/tview"
)

type view struct {
	username string
	password string
	app      *tview.Application
	conn     net.Conn
}

func initView() *view {
	conn, err := net.Dial("tcp", "localhost:9000")
	if err != nil {
		log.Fatal(err)
	}
	return &view{
		app:  tview.NewApplication(),
		conn: conn,
	}
}

func (v *view) genTextArea() (*tview.Flex, *tview.TextArea, *tview.TextView) {

	inputArea := tview.NewTextArea().SetPlaceholder("Enter a new message here...")
	inputArea.SetBorder(true).SetTitle(" Write Here ")

	textView := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetWordWrap(true).
		SetChangedFunc(func() {
			v.app.Draw()
		})
	textView.SetBorder(true).SetTitle(" Messages here ")

	sendButton := tview.NewButton(" Send Message ")
	sendButton.SetSelectedFunc(func() {
		fmt.Fprintf(textView, "YOU: \n%s\n", inputArea.GetText())
		sendConn(v.conn, "POST", v.username, v.password, inputArea.GetText(), "", "")
		inputArea.SetText("", true)
	})

	return tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(textView, 0, 3, false).
			AddItem(inputArea, 0, 1, false).
			AddItem(sendButton, 1, 0, false),
		inputArea, textView
}

func (v *view) getLoginForm() *tview.Form {

	loginForm := tview.NewForm().
		AddTextView(" Login Page Info ", "This is the info", 0, 0, true, true).
		AddInputField(" Username: ", "", 16, nil, nil).
		AddPasswordField(" Password: ", "", 16, '*', nil)

	loginForm.AddButton(" Login ", func() {
		v.username = loginForm.GetFormItemByLabel(" Username: ").(*tview.InputField).GetText()
		v.password = loginForm.GetFormItemByLabel(" Password: ").(*tview.InputField).GetText()
		sendConn(v.conn, "AUTH", v.username, v.password, "", "", "")
	}).
		SetBorder(true).
		SetTitle(" Login Page ")

	return loginForm
}

func (v *view) getRoomList() (*tview.Flex, *tview.List) {

	// TODO: need to add an add room button
	// ROOM_ADD and ROOM_DEL
	// login to rooms
	// allow invites to room
	// TODO: Dont need to have room password to delete room

	inputAddRoomName := tview.NewInputField().SetFieldWidth(14)
	inputAddRoomPass := tview.NewInputField().SetFieldWidth(14)
	inputDelRoomName := tview.NewInputField().SetFieldWidth(14)

	roomFormFlex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(tview.NewTextView().SetText(" room name "), 1, 0, false).
		AddItem(inputAddRoomName, 1, 0, false).
		AddItem(tview.NewTextView().SetText(" room password "), 1, 0, false).
		AddItem(inputAddRoomPass, 1, 0, false).
		AddItem(tview.NewButton(" add room ").SetSelectedFunc(func() {
			sendConn(v.conn, "ADD_ROOM", v.username, v.password, "", inputAddRoomName.GetText(), inputAddRoomPass.GetText())
			inputAddRoomName.SetText("")
			inputAddRoomPass.SetText("")
		}), 1, 0, false).
		AddItem(tview.NewTextView().SetText(" room name "), 1, 0, false).
		AddItem(inputDelRoomName, 1, 0, false).
		AddItem(tview.NewButton(" delete room ").SetSelectedFunc(func() {
			sendConn(v.conn, "DEL_ROOM", v.username, v.password, "", inputDelRoomName.GetText(), "")
			inputDelRoomName.SetText("")
		}), 1, 0, false)

	roomList := tview.NewList().
		AddItem("Home", "Home Main Menu", 'h', nil)
	// AddItem("Room 1", "Some explanatory text", 'a', nil).
	// AddItem("Room 2", "Some explanatory text", 'b', nil).
	// AddItem("Quit", "Press to exit", 'q', func() {
	// 	v.app.Stop()
	// })
	roomList.SetBorder(true).SetTitle(" Your Rooms ")

	roomFlex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(roomFormFlex, 8, 0, false).
		AddItem(roomList, 0, 1, false)

	return roomFlex, roomList
}
