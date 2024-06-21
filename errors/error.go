package errors

import (
	"fmt"
	"strings"

	"github.com/owlto-finance/utils-go/log"
)

type BizError struct {
	Code int64                  `json:"code"`
	Msg  string                 `json:"msg"`
	Info map[string]interface{} `json:"info"`
}

func (e *BizError) Error() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("code: %d, msg: %s", e.Code, e.Msg))
	if len(e.Info) > 0 {
		sb.WriteString(" | info: {")
		infoParts := make([]string, 0, len(e.Info))
		for key, value := range e.Info {
			infoParts = append(infoParts, fmt.Sprintf("%s = %v", key, value))
		}
		sb.WriteString(strings.Join(infoParts, ", "))
		sb.WriteString("}")
	}
	return sb.String()
}

func NewBizError(code int64, msg string) *BizError {
	return &BizError{Code: code, Msg: msg}
}

func (e *BizError) GetCode() int64 {
	return e.Code
}

func (e *BizError) GetMsg() string {
	return e.Msg
}

func (e *BizError) GetMsgPtr() *string {
	return &e.Msg
}

func (e *BizError) WithMsg(msg string) *BizError {
	e.Msg = msg
	return e
}

func (e *BizError) WithInfo(kv ...interface{}) *BizError {
	if len(kv)%2 != 0 {
		fmt.Println("Warning: WithInfo requires an even number of arguments, ignoring the last key-value pair.")
		kv = kv[:len(kv)-1] // Ignore the last key if it's unmatched
	}
	var info = make(map[string]interface{})
	for k, v := range e.Info {
		info[k] = v
	}

	for i := 0; i < len(kv); i += 2 {
		key, ok := kv[i].(string)
		if !ok {
			log.Warn("Warning: WithInfo requires string keys, ignoring this key-value pair.")
			continue
		}
		info[key] = kv[i+1]
	}
	e.Info = info
	return e
}
