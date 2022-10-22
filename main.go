package main

import (
	"github.com/slack-go/slack"
)

func main() {

	api := slack.New("xoxb-3517990543552-3498517567635-CX5hwEl01DUYCDXtxSzF40zp")

	api.PostMessage("C03EJ6VUTKL", slack.MsgOptionText("Hello Customer Service Team", false))

}
