package main

import (
	"os"

	"github.com/slack-go/slack"
)

func getEnv(key, def string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return def
}

func main() {

	api := slack.New(getEnv("BOT_TOKEN", "abcxyz"))

	api.PostMessage(getEnv("CHANNEL_ID", "abcxyz"), slack.MsgOptionText("Hello Customer Service Team", false))

}
