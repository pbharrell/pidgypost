package model

import "github.com/charmbracelet/bubbles/viewport"

func initialViewport() viewport.Model {
	vp := viewport.New(30, 5)
	vp.SetContent(`Welcome to the chat room!
Type a message and press Enter to send.`)
	return vp
}
