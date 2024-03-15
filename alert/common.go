package alert

import (
	"log"
	"sync"
	"time"
)

type Alerter interface {
	AlertText(msg string, err error)
	AlertTextLazy(msg string, err error)
	AlertTextLazyGroup(group string, msg string, err error)
}

type CommonAlerter struct {
	lazyInterval int64
	lazyReset    int64
	lazyTimer    map[string]int64
	mutex        *sync.RWMutex
}

func NewCommonAlerter(lazyInterval int64, lazyReset int64) *CommonAlerter {
	return &CommonAlerter{
		lazyInterval: lazyInterval,
		lazyReset:    lazyReset,
		lazyTimer:    make(map[string]int64),
		mutex:        &sync.RWMutex{},
	}
}

func (ca *CommonAlerter) AlertText(msg string, err error) {
	log.Printf("%s : %v", msg, err)
}

func (ca *CommonAlerter) AlertTextLazyGroup(group string, msg string, err error) {
	ca.DoAlertTextLazy(ca, group, msg, err)
}

func (ca *CommonAlerter) AlertTextLazy(msg string, err error) {
	ca.DoAlertTextLazy(ca, "", msg, err)
}

func (ca *CommonAlerter) DoAlertTextLazy(target Alerter, group string, msg string, err error) {
	ca.mutex.Lock()
	defer ca.mutex.Unlock()
	timer, ok := ca.lazyTimer[group]
	now := time.Now().Unix()
	if !ok {
		ca.lazyTimer[group] = now
	} else {
		if now-timer >= ca.lazyReset {
			ca.lazyTimer[group] = now
			return
		}
		if now-timer >= ca.lazyInterval {
			target.AlertText(msg, err)
			ca.lazyTimer[group] = now
		}
	}
}
