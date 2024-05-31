package utils

import (
	"fmt"
	"log"
)

func Panic(err error, msg interface{}) {
	if err != nil {
		errMsg := fmt.Errorf("Message: %s\nError: %s", msg, err)
		log.Fatal(errMsg)
	}
}
