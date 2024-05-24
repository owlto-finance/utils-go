package alert

import (
	"github.com/go-lark/lark"
)

var LarkBot *Bot

type Bot struct {
	*lark.Bot
}

func NewLarkBot(appID, appSecret string) (*Bot, error) {
	bot := lark.NewChatBot(appID, appSecret)
	err := bot.StartHeartbeat()
	if err != nil {
		return nil, err
	}
	larkBot := &Bot{
		Bot: bot,
	}
	LarkBot = larkBot
	return larkBot, nil
}
