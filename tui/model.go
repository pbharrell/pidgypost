package model

import (
	"fmt"
	"os"
	"strconv"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const gap = "\n\n"

type (
	errMsg error
)

var (
	emptyContact = contact{}

	bubbleStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			Padding(0, 1).
			Margin(1, 0)
)

type MessageType int

const (
	Outgoing MessageType = iota
	Incoming
)

type cachedMessage struct {
	Content     string
	MessageType MessageType
}

type model struct {
	list         list.Model
	listKeys     *listKeyMap
	delegateKeys *listItemDelegateKeyMap
	chatKeys     *chatKeyMap

	viewport    viewport.Model
	messages    []cachedMessage
	textarea    textarea.Model
	senderStyle lipgloss.Style
	err         error

	selected contact
}

func initialModel() model {
	var (
		delegateKeys = newListItemDelegateKeyMap()
		listKeys     = newListKeyMap()
		chatKeys     = newChatKeyMap()
	)

	// FIXME: Load contacts from database
	const numItems = 4
	items := make([]list.Item, numItems)
	for i := range numItems {
		items[i] = contact{"Name" + strconv.FormatInt(int64(i), 10), "Message preview..."}
	}

	// Setup list
	delegate := newListItemDelegate(delegateKeys)
	contactList := list.New(items, delegate, 0, 0)
	contactList.SetShowStatusBar(false)
	contactList.SetShowTitle(false)
	contactList.Title = "Contacts"
	contactList.Styles.Title = titleStyle
	contactList.AdditionalFullHelpKeys = func() []key.Binding {
		return []key.Binding{
			listKeys.insertItem,
			listKeys.clearSelection,
			listKeys.togglePagination,
			listKeys.toggleHelpMenu,
		}
	}

	return model{
		list:         contactList,
		listKeys:     listKeys,
		delegateKeys: delegateKeys,
		chatKeys:     chatKeys,
		textarea:     initialTextArea(),
		messages:     []cachedMessage{},
		viewport:     initialViewport(),
		senderStyle:  lipgloss.NewStyle().Foreground(lipgloss.Color("5")),
		err:          nil,
	}
}

func (m model) Init() tea.Cmd {
	return textarea.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		listWidth := msg.Width / 3
		chatWidth := msg.Width - listWidth

		h, v := appStyle.GetFrameSize()
		m.list.SetSize(listWidth-h, msg.Height-v)

		m.viewport.Width = chatWidth
		m.textarea.SetWidth(chatWidth)
		m.viewport.Height = msg.Height - m.textarea.Height() - 1 // - lipgloss.Height(gap)

		// if len(m.messages) > 0 {
		// Wrap content before setting it.
		// m.viewport.SetContent(lipgloss.NewStyle().Width(m.viewport.Width).Render(strings.Join(m.messages, "\n")))
		// }
		m.viewport.GotoBottom()

	case contactChosenMsg:
		return m.Select(msg.Contact), nil

	case tea.KeyMsg:
		// Don't match any of the keys below if we're actively filtering.
		if m.list.FilterState() == list.Filtering {
			break
		}

		selected := m.Selected()
		if selected == emptyContact {
			// Prevent sending the keymsgs to the chat window components if no contact selected
			switch {

			case key.Matches(msg, m.listKeys.togglePagination):
				m.list.SetShowPagination(!m.list.ShowPagination())
				return m, nil

			case key.Matches(msg, m.listKeys.toggleHelpMenu):
				m.list.SetShowHelp(!m.list.ShowHelp())
				return m, nil

			case key.Matches(msg, m.listKeys.insertItem):
				m.delegateKeys.remove.SetEnabled(true)
				newItem := contact{"New Item!", "New description..."}
				insCmd := m.list.InsertItem(0, newItem)
				statusCmd := m.list.NewStatusMessage(statusMessageStyle("Added " + newItem.Title()))
				return m, tea.Batch(insCmd, statusCmd)

			}

			newListModel, cmd := m.list.Update(msg)
			// This will also call our delegate's update function.
			m.list = newListModel
			cmds = append(cmds, cmd)

		} else {
			switch {
			case key.Matches(msg, m.listKeys.clearSelection):
				return m.ClearSelection(), nil

			case key.Matches(msg, m.chatKeys.send):
				return m.SendMessage(m.textarea.Value()), nil
			}

			newViewport, vpCmd := m.viewport.Update(msg)
			m.viewport = newViewport

			newTextarea, taCmd := m.textarea.Update(msg)
			m.textarea = newTextarea
			cmds = append(cmds, vpCmd, taCmd)
		}
	}

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	// return fmt.Sprintf(
	// 	"%s%s%s",
	// 	m.viewport.View(),
	// 	gap,
	// 	m.textarea.View(),
	// )

	// Get total width, divide it evenly
	listWidth := m.list.Width()
	chatWidth := m.viewport.Width

	// Apply fixed widths
	listStyle := lipgloss.NewStyle().Width(listWidth)
	chatStyle := lipgloss.NewStyle().Width(chatWidth)

	// Chat panel: viewport (top) + textarea (bottom)
	chatPanel := lipgloss.JoinVertical(
		lipgloss.Top,
		m.viewport.View(),
		// gap,
		m.textarea.View(),
	)

	ui := lipgloss.JoinHorizontal(
		lipgloss.Top,
		listStyle.Render(m.list.View()),
		chatStyle.Render(chatPanel),
	)

	// selected := m.Selected()
	// if selected == emptyContact {
	// 	return lipgloss.JoinVertical(lipgloss.Top, appStyle.Render(ui), lipgloss.NewStyle().Foreground(lipgloss.Color("#888")).Render("No selection :("))
	// } else {
	// 	return lipgloss.JoinVertical(lipgloss.Top, appStyle.Render(ui), selected.Title())
	// }

	// return lipgloss.JoinVertical(lipgloss.Top, appStyle.Render(ui), lipgloss.NewStyle().Foreground(lipgloss.Color("#888")).Render("No selection :("))
	return lipgloss.JoinVertical(lipgloss.Top, appStyle.Render(ui))
}

func (m model) Select(c contact) model {
	m.selected = c
	return m
}

func (m model) ClearSelection() model {
	m.selected = emptyContact
	return m
}

func (m model) Selected() contact {
	return m.selected
}

func (m model) SendMessage(message string) model {
	// fmt.Printf(message)
	m.messages = append(m.messages, cachedMessage{Content: message, MessageType: Outgoing})

	var formattedMessages string
	for i := range m.messages {
		formattedMessages += DrawMessage(m.messages[i], m.viewport.Width)
	}

	if len(m.messages) > 0 {
		// Wrap content before setting it.
		m.viewport.SetContent(lipgloss.NewStyle().Width(m.viewport.Width).Render(formattedMessages))
	}

	m.textarea.SetValue("")

	return m
}

func DrawMessage(message cachedMessage, viewportWidth int) string {
	var position lipgloss.Position
	var color string
	switch message.MessageType {
	case Outgoing:
		position = lipgloss.Left
		color = "28"
	case Incoming:
		position = lipgloss.Right
		color = "255"
	}

	bubble := bubbleStyle.BorderForeground(lipgloss.Color(color)).Width(min(lipgloss.Width(message.Content)+2, viewportWidth*3/4)).Render(message.Content)

	return lipgloss.PlaceHorizontal(viewportWidth, position, bubble)
}

func Start() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
