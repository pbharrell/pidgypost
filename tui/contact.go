package model

type contact struct {
	title       string
	description string
}

func (c contact) Title() string       { return c.title }
func (c contact) Description() string { return c.description }
func (c contact) FilterValue() string { return c.title }
