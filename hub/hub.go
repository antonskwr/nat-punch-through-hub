package hub

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
	"sync"
	"time"

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

func handleTCPConnection(conn *net.TCPConn) {
	for {
		// NOTE (antonskwr): in next line hub waits for client to write something to connection
		data, err := bufio.NewReader(conn).ReadString('\n') // NOTE (antonskwr): blocking
		if err != nil {
			if err == io.EOF {
				fmt.Printf("Peer at %s disconnected\n", conn.RemoteAddr().String())
				util.PrintSeparator()
				break
			}
			util.HandleErr(err)
			continue
		}

		if strings.TrimSpace(string(data)) == "STOP" {
			fmt.Printf("Closing connection for client at %s\n", conn.RemoteAddr().String())
			util.PrintSeparator()
			break
		}

		fmt.Print("-> ", string(data))
		t := time.Now()
		hubTime := t.Format(time.RFC3339) + "\n"

		conn.Write([]byte(hubTime))
	}

	conn.Close()
}

func (h *Hub) ListenTCP(port int) error {
	tcpAddr := net.TCPAddr{}
	tcpAddr.Port = port

	tcpListener, err := net.ListenTCP("tcp", &tcpAddr)

	if err != nil {
		return err
	}

	defer tcpListener.Close()

	hubAddr := tcpAddr.String()
	log.Printf("Started TCP hub on %s\n", hubAddr)

	for {
		conn, connErr := tcpListener.AcceptTCP()

		if connErr != nil {
			util.HandleErr(connErr)
			continue
		}

		hAddr := conn.LocalAddr().String()
		rAddr := conn.RemoteAddr().String()

		fmt.Printf("HUB: New connection.\nhub address: %s\nclient address: %s\n", hAddr, rAddr)
		util.PrintSeparator()

		go handleTCPConnection(conn)
	}

	// TODO(antonskwr): currently unreachable
	return nil
}
