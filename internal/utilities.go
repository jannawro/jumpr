package internal

import "log"

func check(msg string, err error) {
	if err != nil {
		log.Fatal(msg, "\n", err)
	}
}
