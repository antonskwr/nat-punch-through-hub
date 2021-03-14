package hub

import (
	"fmt"
	"log"
	"net"
	"sync"

	"github.com/antonskwr/nat-punch-through-hub/util"
)

type Hub struct {
	listLock sync.RWMutex
	gameServers map[uint32]string
}

func NewHub() *Hub {
	h := &Hub{
		gameServers: make(map[uint32]string),
	}
	return h
}

func (h *Hub) ListenTCP(port int) error {
	// TODO(antonskwr): change for secured TLS connection (ListenAndServeTLS)

	tcpAddr := net.TCPAddr{}
	tcpAddr.Port = port

	tcpListener, err := net.ListenTCP("tcp", &tcpAddr)

	if err != nil {
		return err
	}

	hubAddr := tcpAddr.String()
	log.Printf("Started TCP hub on %s\n", hubAddr)

	for {
		conn, connErr := tcpListener.AcceptTCP()

		if connErr != nil {
			util.HandleErr(connErr)
			continue
		}

		sAddr := conn.LocalAddr().String()
		rAddr := conn.RemoteAddr().String()

		fmt.Printf("New connection:\nserver address: %s\nremote address: %s\n", sAddr, rAddr)
		fmt.Println("===============")
	}

	// TODO(antonskwr): currently unreachable
	tcpListener.Close()

	return nil
}
