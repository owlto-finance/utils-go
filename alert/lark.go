package alert

import (
	"fmt"
	"log"

	"github.com/owlto-finance/utils-go/network"
)

// alerter := alert.NewLarkAlerter("https://open.larksuite.com/open-apis/bot/v2/hook/29507972-0dde-4a03-971a-acd012d1d888")
// alerter.AlertText("fuck", errors.New("shit"))

type LarkAlerter struct {
	webhook string
}

func NewLarkAlerter(webhook string) *LarkAlerter {
	return &LarkAlerter{
		webhook: webhook,
	}
}

func (alerter *LarkAlerter) AlertText(msg string, err error) {
	data := map[string]interface{}{
		"msg_type": "text",
		"content": map[string]interface{}{
			"text": fmt.Sprintf("%s : %s", msg, err),
		},
	}
	network.Request(alerter.webhook, data, nil)
	log.Println(msg, ":", err)
}
