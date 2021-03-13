package main

import (
	"log"
	"os"

	"github.com/antonskwr/nat-punch-through-hub/hub"
)

func main() {
	h := hub.NewHub()
	port := getPort()
	log.Printf("Starting hub on port%s ...\n", port)
	log.Fatal(h.ListenAndServe(port))
}

func getPort() string {
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "8080"
	}

	return ":" + port
}
