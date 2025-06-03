package main

import (
	db "github.com/pbharrell/pidgypost/db"
	msg "github.com/pbharrell/pidgypost/msg"
	tui "github.com/pbharrell/pidgypost/tui"
)

func main() {
	msg.NewSentMessage("SENT MESSAGE")
	msg.NewReceivedMessage("RECEIVED MESSAGE")

	tui.Start()
}
