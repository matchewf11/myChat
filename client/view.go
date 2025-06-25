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
		sendConn(v.conn, "POST", v.username, v.password, inputArea.GetText())
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
		sendConn(v.conn, "AUTH", v.username, v.password, "")
	}).
		SetBorder(true).
		SetTitle(" Login Page ")

	return loginForm
}

func (v *view) getRoomList() *tview.List {
	roomList := tview.NewList().
		AddItem("Room 1", "Some explanatory text", 'a', nil).
		AddItem("Room 2", "Some explanatory text", 'b', nil).
		AddItem("Quit", "Press to exit", 'q', func() {
			v.app.Stop()
		})
	roomList.SetBorder(true).SetTitle(" Your Rooms ")
	return roomList
}
