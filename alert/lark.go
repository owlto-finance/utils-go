package alert

import (
	"fmt"
	"log"
	"strings"

	"github.com/owlto-finance/utils-go/network"
)

// alerter := alert.NewLarkAlerter("https://open.larksuite.com/open-apis/bot/v2/hook/29507972-0dde-4a03-971a-acd012d1d888")
// alerter.AlertText("fuck", errors.New("shit"))

type LarkAlerter struct {
	*CommonAlerter
	webhook string
}

func NewLarkAlerter(webhook string) *LarkAlerter {
	return &LarkAlerter{
		webhook:       strings.TrimSpace(webhook),
		CommonAlerter: NewCommonAlerter(180, 900),
	}
}

func (la *LarkAlerter) AlertTextLazyGroup(group string, msg string, err error) {
	la.DoAlertTextLazy(la, group, msg, err)
}

func (la *LarkAlerter) AlertTextLazy(msg string, err error) {
	la.DoAlertTextLazy(la, "", msg, err)
}

func (la *LarkAlerter) AlertText(msg string, err error) {
	data := map[string]interface{}{
		"msg_type": "text",
		"content": map[string]interface{}{
			"text": fmt.Sprintf("%s : %v", msg, err),
		},
	}
	network.Request(la.webhook, data, nil)
	log.Println(msg, ":", err)
}
