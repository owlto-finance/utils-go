package system

import (
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

type QuitCode struct {
	Code int
}

func (e QuitCode) Error() string {
	return strconv.Itoa(e.Code)
}

func WaitForQuitSignals() QuitCode {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	sig := <-sigs
	return QuitCode{Code: int(sig.(syscall.Signal)) + 128}
}
