package model

import (
	"io"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type ContactListItem struct {
	height  int
	spacing int
}

// Render renders the item's view.
func (ContactListItem) Render(w io.Writer, m list.Model, index int, item list.Item) {
	return
}

// Height is the height of the list item.
func (i ContactListItem) Height() int {
	return i.height
}

// Spacing is the size of the horizontal gap between list items in cells.
func (i ContactListItem) Spacing() int {
	return i.spacing
}

// Update is the update loop for items. All messages in the list's update
// loop will pass through here except when the user is setting a filter.
// Use this method to perform item-level updates appropriate to this
// delegate.
func (ContactListItem) Update(msg tea.Msg, m *list.Model) tea.Cmd {
	return nil
}
