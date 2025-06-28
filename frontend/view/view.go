package view

import (
	"encoding/json"
	"log"
	"net"

	"github.com/rivo/tview"
)

type view struct {
	username string
	password string
	App      *tview.Application
	Conn     net.Conn
}

func InitView() *view {
	conn, err := net.Dial("tcp", "localhost:9000")
	if err != nil {
		log.Fatal(err)
	}
	return &view{
		App:  tview.NewApplication(),
		Conn: conn,
	}
}

func (v *view) GetLoginForm() *tview.Form {
	loginForm := tview.NewForm().
		AddTextView(" Login Page Info ", " Updates Here", 0, 0, true, true).
		AddInputField(" Username: ", "", 16, nil, nil).
		AddPasswordField(" Password: ", "", 16, '*', nil)

	loginForm.AddButton(" Login ", func() {
		v.username = loginForm.GetFormItemByLabel(" Username: ").(*tview.InputField).GetText()
		v.password = loginForm.GetFormItemByLabel(" Password: ").(*tview.InputField).GetText()
		sendConn(v.Conn, "LOGIN", v.username, v.password, "", "")
	}).
		SetBorder(true).
		SetTitle(" Login Page ")

	return loginForm
}

func sendConn(conn net.Conn, method, username, password, body, roomName string) {
	if err := json.NewEncoder(conn).Encode(map[string]string{
		"method":    method,
		"username":  username,
		"password":  password,
		"body":      body,
		"room_name": roomName,
	}); err != nil {
		log.Fatal(err)
	}
}

func (v *view) GetRoomList() (*tview.Flex, *tview.List) {
	inputRoomName := tview.NewInputField().SetFieldWidth(14)

	roomFormFlex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(tview.NewTextView().SetText(" room name: "), 1, 0, false).
		AddItem(inputRoomName, 1, 0, false).
		AddItem(tview.NewButton(" add room ").SetSelectedFunc(func() {
			sendConn(v.Conn, "ROOM", v.username, v.password, "", inputRoomName.GetText())
			inputRoomName.SetText("")
		}), 1, 0, false)

	roomList := tview.NewList().
		AddItem("Home", "Home Main Menu", 'h', nil)
	roomList.SetBorder(true).SetTitle(" Your Rooms ")

	roomFlex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(roomFormFlex, 3, 0, false).
		AddItem(roomList, 0, 1, false)

	return roomFlex, roomList
}

func (v *view) GenTextArea() (*tview.Flex, *tview.TextArea, *tview.TextView) {

	inputArea := tview.NewTextArea().SetPlaceholder("Enter a new message here...")
	inputArea.SetBorder(true).SetTitle(" Write Here ")

	textView := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetWordWrap(true).
		SetChangedFunc(func() {
			v.App.Draw()
		})
	textView.SetBorder(true).SetTitle(" Messages here ")

	sendButton := tview.NewButton(" Send Message ")
	sendButton.SetSelectedFunc(func() {
		sendConn(v.Conn, "POST", v.username, v.password, inputArea.GetText(), "")
		inputArea.SetText("", true)
	})

	return tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(textView, 0, 3, false).
			AddItem(inputArea, 0, 1, false).
			AddItem(sendButton, 1, 0, false),
		inputArea, textView
}
