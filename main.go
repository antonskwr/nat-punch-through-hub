package main

import (
	"log"

	"github.com/antonskwr/nat-punch-through-hub/hub"
	"github.com/antonskwr/nat-punch-through-hub/util"
)

func main() {
	h := hub.NewHub()
	port := util.GetPort()
	// log.Printf("Starting TCP hub on port :%d ...\n", port)
	log.Printf("Starting UDP hub on port :%d ...\n", port)
	log.Fatal(h.ListenUDP(port))
}
