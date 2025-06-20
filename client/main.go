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

func main() {

	var (
		username string
		password string
	)

	app := tview.NewApplication()

	conn, err := net.Dial("tcp", "localhost:9000")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	// LOGIN PAGE
	loginPage := tview.NewForm().
		AddInputField(" Username: ", "", 16, nil, nil).
		AddPasswordField(" Password: ", "", 16, '*', nil)
	loginPage.
		AddButton(" Login ", func() {
			username = loginPage.GetFormItemByLabel(" Username: ").(*tview.InputField).GetText()
			password = loginPage.GetFormItemByLabel(" Password: ").(*tview.InputField).GetText()
		}).
		SetBorder(true).
		SetTitle(" Login Page ")

	// TEXT AREA
	inputArea := tview.NewTextArea().SetPlaceholder("Enter a new message here...")
	inputArea.SetBorder(true).SetTitle(" Write Here ")

	// ROOM LIST
	roomList := tview.NewList().
		AddItem("Room 1", "Some explanatory text", 'a', nil).
		AddItem("Room 2", "Some explanatory text", 'b', nil).
		AddItem("Room 3", "Some explanatory text", 'c', nil).
		AddItem("Room 4", "Some explanatory text", 'd', nil).
		AddItem("Quit", "Press to exit", 'q', func() {
			app.Stop()
		})
	roomList.SetBorder(true).SetTitle(" Your Rooms ")

	// TEXT VIEW
	textView := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetWordWrap(true).
		SetChangedFunc(func() {
			app.Draw()
		})
	textView.SetBorder(true).SetTitle(" Messages here ")

	sendButton := tview.NewButton(" Send Message ")
	sendButton.SetSelectedFunc(func() {
		fmt.Fprintf(textView, "YOU: \n%s\n", inputArea.GetText())

		if err := json.NewEncoder(conn).Encode(map[string]string{
			"method":   "POST",
			"body":     inputArea.GetText(),
			"username": username,
			"password": password,
		}); err != nil {
			log.Fatal(err)
		}

		inputArea.SetText("", true)
	})

	// DISPLAY ITEMS
	textArea := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(textView, 0, 3, false).
		AddItem(inputArea, 0, 1, false).
		AddItem(sendButton, 1, 0, false)

	rowFlex := tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(loginPage, 0, 1, true).
		AddItem(roomList, 0, 1, false).
		AddItem(textArea, 0, 1, false)

	reader := bufio.NewReader(conn)

	// start reading for incoming input
	go func() {
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

			if incomingPost.Status != "recieved" {
				continue
			}

			// do somethign with incomingPost
			fmt.Fprintf(textView, "%s: %s\n%s\n", incomingPost.Username, incomingPost.Date, incomingPost.Body)
		}
	}()

	if err := app.SetRoot(rowFlex, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}

}
