package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"

	"github.com/rivo/tview"
)

// add a text view for instructions and debug Info
// make a lua config

type view struct {
	username string
	password string
	app      *tview.Application
	conn     net.Conn
}

func main() {
	view := initView()
	defer view.conn.Close()

	loginPage := view.getLoginPage()
	roomList := view.getRoomList()
	textFlex, inputArea, textView := view.genTextArea()

	rowFlex := tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(loginPage, 0, 1, true).
		AddItem(roomList, 0, 1, false).
		AddItem(textFlex, 0, 1, false)

	go view.readConn(inputArea, textView)

	if err := view.app.SetRoot(rowFlex, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
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

		if err := json.NewEncoder(v.conn).Encode(map[string]string{
			"method":   "POST",
			"body":     inputArea.GetText(),
			"username": v.username,
			"password": v.password,
		}); err != nil {
			log.Fatal(err)
		}

		inputArea.SetText("", true)
	})

	return tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(textView, 0, 3, false).
			AddItem(inputArea, 0, 1, false).
			AddItem(sendButton, 1, 0, false),
		inputArea, textView
}

func (v *view) getLoginPage() *tview.Form {
	loginPage := tview.NewForm().
		AddInputField(" Username: ", "", 16, nil, nil).
		AddPasswordField(" Password: ", "", 16, '*', nil)
	loginPage.
		AddButton(" Login ", func() {

			v.username = loginPage.GetFormItemByLabel(" Username: ").(*tview.InputField).GetText()
			v.password = loginPage.GetFormItemByLabel(" Password: ").(*tview.InputField).GetText()

			if err := json.NewEncoder(v.conn).Encode(map[string]string{
				"method":   "AUTH",
				"username": v.username,
				"password": v.password,
			}); err != nil {
				log.Fatal(err)
			}
		}).
		SetBorder(true).
		SetTitle(" Login Page ")
	return loginPage
}

func (v *view) getRoomList() *tview.List {
	roomList := tview.NewList().
		AddItem("Room 1", "Some explanatory text", 'a', nil).
		AddItem("Room 2", "Some explanatory text", 'b', nil).
		AddItem("Room 3", "Some explanatory text", 'c', nil).
		AddItem("Room 4", "Some explanatory text", 'd', nil).
		AddItem("Quit", "Press to exit", 'q', func() {
			v.app.Stop()
		})
	roomList.SetBorder(true).SetTitle(" Your Rooms ")
	return roomList
}

func (v *view) readConn(inputArea *tview.TextArea, textView *tview.TextView) {

	reader := bufio.NewReader(v.conn)

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err) // fix err handling
		}

		var incomingPost struct {
			Username string `json:"username"`
			Body     string `json:"body"`
			Date     string `json:"date"`
			Status   string `json:"status"`
		}

		if err := json.Unmarshal([]byte(line), &incomingPost); err != nil {
			log.Fatal(err)
		}

		if incomingPost.Status == "loggedin" {
			v.app.SetFocus(inputArea)
			continue
		}

		if incomingPost.Status != "recieved" {
			continue
		}

		// do somethign with incomingPost
		fmt.Fprintf(textView, "%s: %s\n%s\n", incomingPost.Username, incomingPost.Date, incomingPost.Body)
	}

}
