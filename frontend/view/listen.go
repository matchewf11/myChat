package view

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"

	"github.com/rivo/tview"
)

type post struct {
	Date   string `json:"date"`
	Author string `json:"author"`
	Body   string `json:"body"`
}

type room struct {
	Name string `json:"room_name"`
}

func ListenServer(
	view *view,
	loginForm *tview.Form,
	textView *tview.TextView,
	inputArea *tview.TextArea,
	rowFlex *tview.Flex,
	roomList *tview.List,
	textFlex *tview.Flex,
	roomFlex *tview.Flex,
) {

	loginInfo := loginForm.GetFormItemByLabel(" Login Page Info ").(*tview.TextView)

	reader := bufio.NewReader(view.Conn)

	for {

		line, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}

		var update struct {
			Error  string `json:"error"`
			Status string `json:"status"`
			Post   post   `json:"post"`
			Room   room   `json:"room"`
			Posts  []post `json:"posts"`
			Rooms  []room `json:"rooms"`
		}

		if err := json.Unmarshal([]byte(line), &update); err != nil {
			log.Fatal(err)
		}

		if update.Error != "" {
			// write err somewhere
			// on login an main page
			continue
		}

		switch update.Status {
		case "signed up", "logged in":
			loginInfo.SetText(update.Status)
			view.App.QueueUpdateDraw(func() {
				rowFlex.ResizeItem(loginForm, 0, 0)
				rowFlex.ResizeItem(roomFlex, 0, 1)
				rowFlex.ResizeItem(textFlex, 0, 3)
				view.App.SetFocus(inputArea)
			})

			textView.SetText(update.Status + "\n")

			for _, post := range update.Posts {
				fmt.Fprintf(textView, "%s: %s\n%s\n", post.Author, post.Date, post.Body)
			}

			for _, room := range update.Rooms {
				roomList.AddItem(room.Name, "password will go here", 'c', nil)
			}

		case "new_post":
			fmt.Fprintf(textView, "%s: %s\n%s\n", update.Post.Author, update.Post.Date, update.Post.Body)
		case "new_room":
			view.App.QueueUpdateDraw(func() {
				roomList.AddItem(update.Room.Name, "password will go here", 'c', nil)
			})
		default:
			// write an err somewhere
		}
	}
}
