package main

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func main() {
	app := tview.NewApplication()

	var layout *tview.Flex
	var box *tview.Flex
	var form *tview.Form

	// Chat display (read-only text view)
	chatView := tview.NewTextView().
		SetDynamicColors(true).
		SetChangedFunc(func() { app.Draw() })
	chatView.SetBorder(true).SetTitle(" Chat ")

	// Chat input (single line input field)
	chatInput := tview.NewInputField().
		SetLabel("Message: ").
		SetFieldBackgroundColor(tcell.ColorBlack).
		SetFieldTextColor(tcell.ColorWhite)

	// Send message on Enter
	chatInput.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			msg := chatInput.GetText()
			if len(msg) > 0 {
				fmt.Fprintf(chatView, "[yellow]You:[-] %s\n", msg)
				chatView.ScrollToEnd() // scroll to bottom
				chatInput.SetText("")
			}
		}
	})

	// Box holds chatView + chatInput vertically
	box = tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(chatView, 0, 3, false). // chatView takes 3 parts
		AddItem(chatInput, 3, 1, true)  // chatInput 3 lines height

	box.SetBorder(true).
		SetBorderAttributes(tcell.AttrBold).
		SetTitle("[green:black:bu] Welcome to My Chat! ").
		SetTitleAlign(tview.AlignCenter)

	form = tview.NewForm().
		AddInputField("Username", "", 20, nil, nil).
		AddPasswordField("Password", "", 20, '*', nil).
		AddButton("Login", func() {
			layout.ResizeItem(box, 0, 1)  // box expands fully
			layout.ResizeItem(form, 0, 0) // form shrinks

			app.SetFocus(chatInput)
		}).
		AddButton("Quit", func() {
			app.Stop()
		})

	form.
		SetLabelColor(tcell.ColorLimeGreen).
		SetFieldTextColor(tcell.ColorCoral).
		SetFieldBackgroundColor(tcell.ColorBlack).
		SetButtonTextColor(tcell.ColorLimeGreen).
		SetButtonBackgroundColor(tcell.ColorBlack).
		SetBorder(true).
		SetTitle(" Login ").
		SetTitleColor(tcell.ColorLimeGreen).
		SetTitleAlign(tview.AlignLeft)

	layout = tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(box, 0, 0, false).
		AddItem(form, 0, 1, true)

	if err := app.SetRoot(layout, true).EnableMouse(true).EnablePaste(true).Run(); err != nil {
		panic(err)
	}
}

// func main() {
//
// 	conn, err := net.Dial("tcp", "localhost:9000")
// 	if err != nil {
// 		log.Fatal("could not connect to server")
// 	}
// 	defer conn.Close()
//
// 	p := tea.NewProgram(initialModel(conn))
// 	if _, err := p.Run(); err != nil {
// 		log.Fatal(err)
// 	}
//
// }
//
// type model struct {
// 	conn        net.Conn
// 	viewport    viewport.Model
// 	messages    []string
// 	textarea    textarea.Model
// 	senderStyle lipgloss.Style
// 	err         error
// }
//
// type errMsg error
//
// type postMsg struct {
// 	Username string `json:"username"`
// 	Body     string `json:"body"`
// 	Date     string `json:"date"`
// }
//
// const gap = "\n\n"
//
// func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
// 	var (
// 		tiCmd tea.Cmd
// 		vpCmd tea.Cmd
// 	)
//
// 	m.textarea, tiCmd = m.textarea.Update(msg)
// 	m.viewport, vpCmd = m.viewport.Update(msg)
//
// 	switch msg := msg.(type) {
// 	// I added here
// 	case postMsg:
// 		//m.messages = append(m.messages, m.senderStyle.Render("You: ")+msg.Body)
// 		m.messages = append(m.messages, m.senderStyle.Render(msg.Username+" "+msg.Date+" "+msg.Body))
// 		m.viewport.SetContent(lipgloss.NewStyle().Width(m.viewport.Width).Render(strings.Join(m.messages, "\n")))
// 		m.viewport.GotoBottom()
// 		return m, readConn(m.conn)
// 	// I added above here
// 	case tea.WindowSizeMsg:
// 		m.viewport.Width = msg.Width
// 		m.textarea.SetWidth(msg.Width)
// 		m.viewport.Height = msg.Height - m.textarea.Height() - lipgloss.Height(gap)
//
// 		if len(m.messages) > 0 {
// 			// Wrap content before setting it.
// 			m.viewport.SetContent(lipgloss.NewStyle().Width(m.viewport.Width).Render(strings.Join(m.messages, "\n")))
// 		}
// 		m.viewport.GotoBottom()
// 	case tea.KeyMsg:
// 		switch msg.Type {
// 		case tea.KeyCtrlC, tea.KeyEsc:
// 			fmt.Println(m.textarea.Value())
// 			return m, tea.Quit
// 		case tea.KeyEnter:
//
// 			if err := json.NewEncoder(m.conn).Encode(map[string]string{
// 				"method":   "POST",
// 				"username": "Default user",
// 				"body":     m.textarea.Value(),
// 			}); err != nil {
// 				return m, tea.Quit // err later
// 			}
//
// 			m.textarea.Reset()
//
// 		}
// 	case errMsg:
// 		m.err = msg
// 		return m, nil
// 	}
//
// 	return m, tea.Batch(tiCmd, vpCmd)
// }
//
// // Initial Cmd for it to run
// func (m model) Init() tea.Cmd {
// 	return tea.Batch(
// 		textarea.Blink,
// 		readConn(m.conn),
// 	)
// }
//
// // Manages how the model will be displayed
// func (m model) View() string {
// 	return fmt.Sprintf(
// 		"%s%s%s",
// 		m.viewport.View(),
// 		gap,
// 		m.textarea.View(),
// 	)
// }
//
// // Gets the model
// func initialModel(conn net.Conn) model {
//
// 	// sets up the text area
// 	ta := textarea.New()
// 	ta.Placeholder = "Send a message..."
// 	ta.Focus()
// 	ta.Prompt = "┃ "
// 	ta.CharLimit = 280
// 	ta.SetWidth(30)
// 	ta.SetHeight(3)
// 	ta.FocusedStyle.CursorLine = lipgloss.NewStyle() // no cursor line styling
// 	ta.ShowLineNumbers = false
// 	ta.KeyMap.InsertNewline.SetEnabled(false)
//
// 	// sets up the place with the messages
// 	vp := viewport.New(30, 5)
// 	vp.SetContent(`Welcome to the chat room!
// Type a message and press Enter to send.`)
//
// 	return model{
// 		textarea:    ta,
// 		viewport:    vp,
// 		messages:    []string{},
// 		senderStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("5")),
// 		conn:        conn,
// 		err:         nil,
// 	}
// }
//
// func readConn(conn net.Conn) tea.Cmd {
// 	return func() tea.Msg {
//
// 		reader := bufio.NewReader(conn)
//
// 		line, err := reader.ReadString('\n')
// 		if err != nil {
// 			return errMsg(err)
// 		}
//
// 		var incomingPost struct {
// 			Username string `json:"username"`
// 			Body     string `json:"body"`
// 			Date     string `json:"date"`
// 			Status   string `json:"status"`
// 		}
//
// 		if err := json.Unmarshal([]byte(line), &incomingPost); err != nil {
// 			return errMsg(err)
// 		}
//
// 		if incomingPost.Status != "recieved" {
// 			return nil
// 		}
//
// 		return postMsg{
// 			Username: incomingPost.Username,
// 			Body:     incomingPost.Body,
// 			Date:     incomingPost.Date,
// 		}
// 	}
// }
