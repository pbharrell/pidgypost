package model

import "github.com/charmbracelet/bubbles/key"

type listKeyMap struct {
	togglePagination key.Binding
	toggleHelpMenu   key.Binding
	insertItem       key.Binding
	clearSelection   key.Binding
}

type listItemDelegateKeyMap struct {
	choose key.Binding
	remove key.Binding
}

type chatKeyMap struct {
	send key.Binding
}

func newListKeyMap() *listKeyMap {
	return &listKeyMap{
		insertItem: key.NewBinding(
			key.WithKeys("a"),
			key.WithHelp("a", "add item"),
		),
		togglePagination: key.NewBinding(
			key.WithKeys("P"),
			key.WithHelp("P", "toggle pagination"),
		),
		toggleHelpMenu: key.NewBinding(
			key.WithKeys("H"),
			key.WithHelp("H", "toggle help"),
		),
		clearSelection: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "clear selection"),
		),
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

func newChatKeyMap() *chatKeyMap {
	return &chatKeyMap{
		send: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "send"),
		),
	}
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
