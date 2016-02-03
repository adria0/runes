package store

import (
	"log"
)

func assert(err error, msg string) {
	if err != nil {
		log.Fatalln(msg, err)
	}
}
