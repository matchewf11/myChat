package main

//conn.Write([]byte("john\n")) -> this is what client should do

// make a display area for in comiing error messages and success messages

// import (
// 	"bufio"
// 	"encoding/json"
// 	"fmt"
// 	"log"
// 	"net"

// 	"github.com/rivo/tview"
// )

// const (
// 	StatusLoginFail = "login fail"
// 	StatusLoggedIn  = "loggedin"
// )

// func main() {

// 	view := initView()
// 	defer view.conn.Close()

// 	loginForm := view.getLoginForm()
// 	roomFlex, roomList := view.getRoomList()
// 	textFlex, inputArea, textView := view.genTextArea()

// 	rowFlex := tview.NewFlex().SetDirection(tview.FlexColumn).
// 		AddItem(loginForm, 0, 1, true).
// 		AddItem(roomFlex, 0, 0, false).
// 		AddItem(textFlex, 0, 0, false)

// 	go listenServer(view, loginForm, textView, inputArea, rowFlex, roomList, textFlex, roomFlex)

// 	if err := view.app.SetRoot(rowFlex, true).EnableMouse(true).Run(); err != nil {
// 		log.Fatal(err)
// 	}
// }

// type view struct {
// 	username string
// 	password string
// 	app      *tview.Application
// 	conn     net.Conn
// }

// func initView() *view {
// 	conn, err := net.Dial("tcp", "localhost:9000")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	return &view{
// 		app:  tview.NewApplication(),
// 		conn: conn,
// 	}
// }

// func (v *view) genTextArea() (*tview.Flex, *tview.TextArea, *tview.TextView) {

// 	inputArea := tview.NewTextArea().SetPlaceholder("Enter a new message here...")
// 	inputArea.SetBorder(true).SetTitle(" Write Here ")

// 	textView := tview.NewTextView().
// 		SetDynamicColors(true).
// 		SetRegions(true).
// 		SetWordWrap(true).
// 		SetChangedFunc(func() {
// 			v.app.Draw()
// 		})
// 	textView.SetBorder(true).SetTitle(" Messages here ")

// 	sendButton := tview.NewButton(" Send Message ")
// 	sendButton.SetSelectedFunc(func() {
// 		sendConn(v.conn, "POST", v.username, v.password, inputArea.GetText(), "", "")
// 		inputArea.SetText("", true)
// 	})

// 	return tview.NewFlex().SetDirection(tview.FlexRow).
// 			AddItem(textView, 0, 3, false).
// 			AddItem(inputArea, 0, 1, false).
// 			AddItem(sendButton, 1, 0, false),
// 		inputArea, textView
// }

// func (v *view) getLoginForm() *tview.Form {

// 	loginForm := tview.NewForm().
// 		AddTextView(" Login Page Info ", "This is the info", 0, 0, true, true).
// 		AddInputField(" Username: ", "", 16, nil, nil).
// 		AddPasswordField(" Password: ", "", 16, '*', nil)

// 	loginForm.AddButton(" Login ", func() {
// 		v.username = loginForm.GetFormItemByLabel(" Username: ").(*tview.InputField).GetText()
// 		v.password = loginForm.GetFormItemByLabel(" Password: ").(*tview.InputField).GetText()
// 		sendConn(v.conn, "AUTH", v.username, v.password, "", "", "")
// 	}).
// 		SetBorder(true).
// 		SetTitle(" Login Page ")

// 	return loginForm
// }

// func (v *view) getRoomList() (*tview.Flex, *tview.List) {

// 	// TODO: need to add an add room button
// 	// ROOM_ADD and ROOM_DEL
// 	// login to rooms
// 	// allow invites to room
// 	// TODO: Dont need to have room password to delete room

// 	inputAddRoomName := tview.NewInputField().SetFieldWidth(14)
// 	inputAddRoomPass := tview.NewInputField().SetFieldWidth(14)
// 	inputDelRoomName := tview.NewInputField().SetFieldWidth(14)

// 	roomFormFlex := tview.NewFlex().SetDirection(tview.FlexRow).
// 		AddItem(tview.NewTextView().SetText(" room name "), 1, 0, false).
// 		AddItem(inputAddRoomName, 1, 0, false).
// 		AddItem(tview.NewTextView().SetText(" room password "), 1, 0, false).
// 		AddItem(inputAddRoomPass, 1, 0, false).
// 		AddItem(tview.NewButton(" add room ").SetSelectedFunc(func() {
// 			sendConn(v.conn, "ADD_ROOM", v.username, v.password, "", inputAddRoomName.GetText(), inputAddRoomPass.GetText())
// 			inputAddRoomName.SetText("")
// 			inputAddRoomPass.SetText("")
// 		}), 1, 0, false).
// 		AddItem(tview.NewTextView().SetText(" room name "), 1, 0, false).
// 		AddItem(inputDelRoomName, 1, 0, false).
// 		AddItem(tview.NewButton(" delete room ").SetSelectedFunc(func() {
// 			sendConn(v.conn, "DEL_ROOM", v.username, v.password, "", inputDelRoomName.GetText(), "")
// 			inputDelRoomName.SetText("")
// 		}), 1, 0, false)

// 	roomList := tview.NewList().
// 		AddItem("Home", "Home Main Menu", 'h', nil)
// 	// AddItem("Room 1", "Some explanatory text", 'a', nil).
// 	// AddItem("Room 2", "Some explanatory text", 'b', nil).
// 	// AddItem("Quit", "Press to exit", 'q', func() {
// 	// 	v.app.Stop()
// 	// })
// 	roomList.SetBorder(true).SetTitle(" Your Rooms ")

// 	roomFlex := tview.NewFlex().
// 		SetDirection(tview.FlexRow).
// 		AddItem(roomFormFlex, 8, 0, false).
// 		AddItem(roomList, 0, 1, false)

// 	return roomFlex, roomList
// }

// func sendConn(conn net.Conn, method, username, password, body, roomName, roomPass string) {
// 	if err := json.NewEncoder(conn).Encode(map[string]string{
// 		"method":    method,
// 		"username":  username,
// 		"password":  password,
// 		"room_name": roomName,
// 		"room_pass": roomPass,
// 		"body":      body,
// 	}); err != nil {
// 		log.Fatal(err)
// 	}
// }

// func listenServer(
// 	view *view,
// 	loginForm *tview.Form,
// 	textView *tview.TextView,
// 	inputArea *tview.TextArea,
// 	rowFlex *tview.Flex,
// 	roomList *tview.List,
// 	textFlex *tview.Flex,
// 	roomFlex *tview.Flex,
// ) {

// 	reader := bufio.NewReader(view.conn)
// 	loginInfo := loginForm.GetFormItemByLabel(" Login Page Info ").(*tview.TextView)
// 	for {

// 		line, err := reader.ReadString('\n')

// 		if err != nil {
// 			log.Fatal(err)
// 		}

// 		type message struct {
// 			Date   string `json:"date"`
// 			Author string `json:"author"`
// 			Body   string `json:"body"`
// 		}

// 		var incomingPost struct {
// 			Username string    `json:"username"`
// 			Body     string    `json:"body"`
// 			Date     string    `json:"date"`
// 			Status   string    `json:"status"`
// 			Messages []message `json:"messages"`
// 		}

// 		if err := json.Unmarshal([]byte(line), &incomingPost); err != nil {
// 			log.Fatal(err)
// 		}

// 		switch incomingPost.Status {
// 		case StatusLoginFail:
// 			loginInfo.SetText(incomingPost.Body)
// 			fmt.Fprintf(textView, "%s: %s\n%s\n", incomingPost.Username, incomingPost.Date, incomingPost.Body)
// 		case StatusLoggedIn:
// 			loginInfo.SetText(incomingPost.Body)
// 			view.app.QueueUpdateDraw(func() {
// 				rowFlex.ResizeItem(loginForm, 0, 0)
// 				rowFlex.ResizeItem(roomFlex, 0, 1)
// 				rowFlex.ResizeItem(textFlex, 0, 3)
// 				view.app.SetFocus(inputArea)
// 			})

// 			if incomingPost.Date != "" {
// 				textView.SetText("Last Login: " + incomingPost.Date + "\n")
// 			} else {
// 				textView.SetText("Welcome First Time User\n")
// 			}

// 			for _, mes := range incomingPost.Messages {
// 				fmt.Fprintf(textView, "%s: %s\n%s\n\n", mes.Author, mes.Date, mes.Body)
// 			}

// 			if incomingPost.Body != "logged in" && incomingPost.Username != "" {
// 				fmt.Fprintf(textView, "%s: %s\n%s\n", incomingPost.Username, incomingPost.Date, incomingPost.Body)
// 			}

// 		case "recieved":
// 			// no special behaviors
// 		default:
// 			// unknown status
// 		}
// 	}
// }
