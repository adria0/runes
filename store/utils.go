package store

import (
	"errors"
	"fmt"
	"log"
)

func forwardErrorAndClose(ch StoreChannel, err error) {
	fmt.Println("ERROR: " + err.Error())
	result := StoreResult{}
	result.Err = err
	ch <- result
	close(ch)
}

func sendErrorAndClose(ch StoreChannel, message string) {
	fmt.Println("ERROR: " + message)
	result := StoreResult{}
	result.Err = errors.New(message)
	ch <- result
	close(ch)
}

func sendSuccessAndClose(ch StoreChannel, data interface{}) {
	result := StoreResult{}
	result.Err = nil
	result.Data = data
	ch <- result
	close(ch)
}

func assert(err error, msg string) {
	if err != nil {
		log.Fatalln(msg, err)
	}
}
