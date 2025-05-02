package model

import (
	"log"
	"os"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	appStyle = lipgloss.NewStyle().Padding(1, 2)

	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFDF5")).
			Background(lipgloss.Color("#25A065")).
			Padding(0, 1)

	statusMessageStyle = lipgloss.NewStyle().
				Foreground(lipgloss.AdaptiveColor{Light: "#04B575", Dark: "#04B575"}).
				Render

	logger *log.Logger
)

type contactChosenMsg struct {
	Contact contact
}

func OverrideListItemStyles(s *list.DefaultItemStyles) {
	s.NormalTitle = lipgloss.NewStyle().
		Foreground(lipgloss.AdaptiveColor{Light: "#1a1a1a", Dark: "#dddddd"}).
		Padding(0, 0, 0, 2) //nolint:mnd

	s.NormalDesc = s.NormalTitle.
		Foreground(lipgloss.AdaptiveColor{Light: "#A49FA5", Dark: "#777777"})

	s.SelectedTitle = lipgloss.NewStyle().
		Border(lipgloss.NormalBorder(), false, false, false, true).
		BorderForeground(lipgloss.AdaptiveColor{Light: "#009933", Dark: "#006622"}).
		Foreground(lipgloss.AdaptiveColor{Light: "#00b33c", Dark: "#00b33c"}).
		Padding(0, 0, 0, 1)

	s.SelectedDesc = s.SelectedTitle.
		Foreground(lipgloss.AdaptiveColor{Light: "#009933", Dark: "#006622"})

	s.DimmedTitle = lipgloss.NewStyle().
		Foreground(lipgloss.AdaptiveColor{Light: "#A49FA5", Dark: "#777777"}).Padding(0, 0, 0, 2)

	s.DimmedDesc = s.DimmedTitle.
		Foreground(lipgloss.AdaptiveColor{Light: "#C2B8C2", Dark: "#4D4D4D"})

	s.FilterMatch = lipgloss.NewStyle().Underline(true)
}

func Update(keys *listItemDelegateKeyMap, msg tea.Msg, m *list.Model) tea.Cmd {
	var title string

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.choose):
			if i, ok := m.SelectedItem().(contact); ok {
				selectedContactMsg := func() tea.Msg { return contactChosenMsg{Contact: i} }
				stateMsg := m.NewStatusMessage(statusMessageStyle("You chose " + i.Title()))
				return tea.Batch(selectedContactMsg, stateMsg)
			}

		case key.Matches(msg, keys.remove):
			logger.Println("Something was removed")
			index := m.Index()
			m.RemoveItem(index)
			if len(m.Items()) == 0 {
				keys.remove.SetEnabled(false)
			}
			return m.NewStatusMessage(statusMessageStyle("Deleted " + title))
		}
	}

	return nil
}

func newListItemDelegate(keys *listItemDelegateKeyMap) list.DefaultDelegate {
	f, err := os.OpenFile("debug.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	logger = log.New(f, "", log.LstdFlags|log.Lshortfile)

	d := list.NewDefaultDelegate()
	OverrideListItemStyles(&d.Styles)

	d.UpdateFunc = func(msg tea.Msg, m *list.Model) tea.Cmd {
		return Update(keys, msg, m)
	}

	help := []key.Binding{keys.choose, keys.remove}

	d.ShortHelpFunc = func() []key.Binding {
		return help
	}

	d.FullHelpFunc = func() [][]key.Binding {
		return [][]key.Binding{help}
	}

	return d
}

type listItemDelegateKeyMap struct {
	choose key.Binding
	remove key.Binding
}

// Additional short help entries. This satisfies the help.KeyMap interface and
// is entirely optional.
func (d listItemDelegateKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		d.choose,
		d.remove,
	}
}

// Additional full help entries. This satisfies the help.KeyMap interface and
// is entirely optional.
func (d listItemDelegateKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{
			d.choose,
			d.remove,
		},
	}
}

func newListItemDelegateKeyMap() *listItemDelegateKeyMap {
	return &listItemDelegateKeyMap{
		choose: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "choose"),
		),
		remove: key.NewBinding(
			key.WithKeys("x", "backspace"),
			key.WithHelp("x", "delete"),
		),
	}
}
