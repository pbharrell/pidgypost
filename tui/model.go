package model

import (
	"fmt"
	"os"
	"strconv"
	"strings"

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

var emptyContact = contact{}

type mode struct {
}

type model struct {
	list         list.Model
	listKeys     *listKeyMap
	delegateKeys *listItemDelegateKeyMap

	viewport    viewport.Model
	messages    []string
	textarea    textarea.Model
	senderStyle lipgloss.Style
	err         error

	selected contact
}

func initialModel() model {
	var (
		delegateKeys = newListItemDelegateKeyMap()
		listKeys     = newListKeyMap()
	)

	// FIXME: Load contacts from database
	const numItems = 4
	items := make([]list.Item, numItems)
	for i := range numItems {
		items[i] = contact{"Name" + strconv.FormatInt(int64(i), 10), "Message preview..."}
	}

	// Setup list
	delegate := newListItemDelegate(delegateKeys)
	groceryList := list.New(items, delegate, 0, 0)
	groceryList.Title = "Groceries"
	groceryList.Styles.Title = titleStyle
	groceryList.AdditionalFullHelpKeys = func() []key.Binding {
		return []key.Binding{
			listKeys.toggleSpinner,
			listKeys.insertItem,
			listKeys.clearSelection,
			listKeys.toggleTitleBar,
			listKeys.toggleStatusBar,
			listKeys.togglePagination,
			listKeys.toggleHelpMenu,
		}
	}

	ta := textarea.New()
	ta.Placeholder = "Send a message..."
	ta.Focus()

	ta.Prompt = "┃ "
	ta.CharLimit = 280

	ta.SetWidth(30)
	ta.SetHeight(3)

	// Remove cursor line styling
	ta.FocusedStyle.CursorLine = lipgloss.NewStyle()

	ta.ShowLineNumbers = false

	vp := viewport.New(30, 5)
	vp.SetContent(`Welcome to the chat room!
Type a message and press Enter to send.`)

	ta.KeyMap.InsertNewline.SetEnabled(false)

	return model{
		list:         groceryList,
		listKeys:     listKeys,
		delegateKeys: delegateKeys,
		textarea:     ta,
		messages:     []string{},
		viewport:     vp,
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
		m.viewport.Height = msg.Height - m.textarea.Height() - lipgloss.Height(gap)

		if len(m.messages) > 0 {
			// Wrap content before setting it.
			m.viewport.SetContent(lipgloss.NewStyle().Width(m.viewport.Width).Render(strings.Join(m.messages, "\n")))
		}
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
			// Prevent sending the keymsgs to the individual components if selected
			logger.Println("Made it to non-selected key handlers. Selected: ", selected.Title())
			switch {
			case key.Matches(msg, m.listKeys.toggleSpinner):
				cmd := m.list.ToggleSpinner()
				return m, cmd

			case key.Matches(msg, m.listKeys.toggleTitleBar):
				v := !m.list.ShowTitle()
				m.list.SetShowTitle(v)
				m.list.SetShowFilter(v)
				m.list.SetFilteringEnabled(v)
				return m, nil

			case key.Matches(msg, m.listKeys.toggleStatusBar):
				m.list.SetShowStatusBar(!m.list.ShowStatusBar())
				return m, nil

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
		} else {
			logger.Println("Made it to selected key handlers. Selected: ", selected.Title())
			switch {
			case key.Matches(msg, m.listKeys.clearSelection):
				return m.ClearSelection(), nil
			}
		}
	}

	// This will also call our delegate's update function.
	newListModel, cmd := m.list.Update(msg)
	m.list = newListModel
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
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
		gap,
		m.textarea.View(),
	)

	ui := lipgloss.JoinHorizontal(
		lipgloss.Top,
		listStyle.Render(m.list.View()),
		chatStyle.Render(chatPanel),
	)

	selected := m.Selected()
	logger.Println("in view, selected: ", selected.Title())
	if selected == emptyContact {
		return lipgloss.JoinVertical(lipgloss.Top, appStyle.Render(ui), lipgloss.NewStyle().Foreground(lipgloss.Color("#888")).Render("No selection :("))
	} else {
		return lipgloss.JoinVertical(lipgloss.Top, appStyle.Render(ui), selected.Title())
	}
}

func (m model) Select(c contact) model {
	m.selected = c
	logger.Println("Contact with name", m.selected.Title(), "selected")
	return m
}

func (m model) ClearSelection() model {
	m.selected = emptyContact
	return m
}

func (m model) Selected() contact {
	return m.selected
}

func Start() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
