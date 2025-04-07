package main

import (
	msg "github.com/pbharrell/pidgypost/msg"
	tui "github.com/pbharrell/pidgypost/tui"
)

func main() {
	msg.NewMessage("MESSAGE")

	tui.Start()
}
