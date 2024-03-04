package alert

type Alerter interface {
	AlertText(msg string, err error)
}
