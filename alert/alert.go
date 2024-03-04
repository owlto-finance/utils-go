package alert

import "log"

func Alert(msg string, err error) {
	log.Println(err)
}
