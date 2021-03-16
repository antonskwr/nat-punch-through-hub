package main

import (
	"log"
	"os"
	"strconv"

	"github.com/antonskwr/nat-punch-through-hub/hub"
)

func main() {
	h := hub.NewHub()
	port := getPort()
	// log.Printf("Starting TCP hub on port :%d ...\n", port)
	log.Printf("Starting UDP hub on port :%d ...\n", port)
	log.Fatal(h.ListenUDP(port))
}

func getPort() int {
	port := os.Getenv("PORT")
	intPort, err := strconv.Atoi(port)

	if err != nil {
		intPort = 8080
	}

	return intPort
}
