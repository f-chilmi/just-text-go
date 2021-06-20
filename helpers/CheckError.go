package helpers

import "log"

func CheckError(msg string, err error) {
	if err != nil {
		log.Fatalf("msg. %v", err)
	}
}
