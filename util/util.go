package util

import (
	"fmt"
	"log"
	"os"
	"strconv"
)

func HandleErrNonFatal(err error, message ...string) {
	if err != nil {
		if len(message) > 0 {
			err = fmt.Errorf("[%s] -- %w --", message[0], err)
		}
		log.Printf("%v", err)
	}
}

func HandleErrFatal(err error, message ...string) {
	if err != nil {
		if len(message) > 0 {
			err = fmt.Errorf("[%s] -- %w --", message[0], err)
		}
		log.Fatal(err)
	}
}

func PrintSeparator() {
	fmt.Printf("=========\n\n")
}

func GetPort() int {
	port := os.Getenv("PORT")
	intPort, err := strconv.Atoi(port)

	if err != nil {
		intPort = 8080
	}

	return intPort
}
