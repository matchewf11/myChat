package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"

	"github.com/rivo/tview"
)

func sendConn(conn net.Conn, method, username, password, body, roomName, roomPass string) {
	if err := json.NewEncoder(conn).Encode(map[string]string{
		"method":    method,
		"username":  username,
		"password":  password,
		"room_name": roomName,
		"room_pass": roomPass,
		"body":      body,
	}); err != nil {
		log.Fatal(err)
	}
}

func listenServer(
	view *view,
	loginForm *tview.Form,
	textView *tview.TextView,
	inputArea *tview.TextArea,
	rowFlex *tview.Flex,
	roomList *tview.List,
	textFlex *tview.Flex,
	roomFlex *tview.Flex,
) {

	reader := bufio.NewReader(view.conn)
	loginInfo := loginForm.GetFormItemByLabel(" Login Page Info ").(*tview.TextView)
	for {

		line, err := reader.ReadString('\n')

		if err != nil {
			log.Fatal(err)
		}

		type message struct {
			Date   string `json:"date"`
			Author string `json:"author"`
			Body   string `json:"body"`
		}

		var incomingPost struct {
			Username string    `json:"username"`
			Body     string    `json:"body"`
			Date     string    `json:"date"`
			Status   string    `json:"status"`
			Messages []message `json:"messages"`
		}

		if err := json.Unmarshal([]byte(line), &incomingPost); err != nil {
			log.Fatal(err)
		}

		switch incomingPost.Status {
		case StatusLoginFail:
			loginInfo.SetText(incomingPost.Body)
			fmt.Fprintf(textView, "%s: %s\n%s\n", incomingPost.Username, incomingPost.Date, incomingPost.Body)
		case StatusLoggedIn:
			loginInfo.SetText(incomingPost.Body)
			view.app.QueueUpdateDraw(func() {
				rowFlex.ResizeItem(loginForm, 0, 0)
				rowFlex.ResizeItem(roomFlex, 0, 1)
				rowFlex.ResizeItem(textFlex, 0, 3)
				view.app.SetFocus(inputArea)
			})

			if incomingPost.Date != "" {
				textView.SetText("Last Login: " + incomingPost.Date + "\n")
			} else {
				textView.SetText("Welcome First Time User\n")
			}

			for _, mes := range incomingPost.Messages {
				fmt.Fprintf(textView, "%s: %s\n%s\n\n", mes.Author, mes.Date, mes.Body)
			}

			if incomingPost.Body != "logged in" && incomingPost.Username != "" {
				fmt.Fprintf(textView, "%s: %s\n%s\n", incomingPost.Username, incomingPost.Date, incomingPost.Body)
			}

		case "recieved":
			// no special behaviors
		default:
			// unknown status
		}
	}
}
