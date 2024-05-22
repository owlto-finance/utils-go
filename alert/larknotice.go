package alert

import (
	"github.com/go-lark/lark"
	"log"
)

var LarkBot *Bot

type Bot struct {
	*lark.Bot
}

func init() {
	bot := lark.NewChatBot("cli_a6c04dc5b7b81010", "PKM6eJPvJOQP79hFtNXgWhlThrXZkKPt")
	err := bot.StartHeartbeat()
	if err != nil {
		log.Fatal("start lark notice heartbeat error:", err)
	}
	LarkBot = &Bot{
		Bot: bot,
	}
}
