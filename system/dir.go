package system

import (
	"os"
)

func MakeDirAll(home string) error {
	_, err := os.Stat(home)
	if err == nil {
		return nil
	} else {
		if !os.IsNotExist(err) {
			return err
		}
	}

	if err := os.MkdirAll(home, 0755); err != nil {
		return err
	}

	return nil
}
