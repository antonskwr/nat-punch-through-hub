package hub

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/antonskwr/nat-punch-through-hub/util"
)

type Hub struct {
	listLock sync.RWMutex
	gameServers map[int]string
}

func NewHub() *Hub {
	h := &Hub{
		gameServers: make(map[int]string),
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
}

func (h *Hub) HandleMsgUDP(msg string, addr *net.UDPAddr) string {
	resp := "Unknown msg"
	splittedMsgs := strings.Split(msg, " ")

	if len(splittedMsgs) == 1 {
		if splittedMsgs[0] == "LIST" {
			if len(h.gameServers) == 0 {
				return "No hosts registered"
			}
			resp = "Listing hosts ...\n"
			for id := range h.gameServers {
				resp = fmt.Sprintf("%sHost[%d]\n", resp, id)
			}
		}
	} else if len(splittedMsgs) == 2 {
		id, err := strconv.Atoi(splittedMsgs[1])
		if err == nil {
			switch splittedMsgs[0] {
			case "ADD":
				h.listLock.Lock()
				h.gameServers[id] = addr.String()
				h.listLock.Unlock()
				resp = "Host added successfully"
			case "CONN":
				// TODO(antonskwr): start NAT punch through
			}
		}
	}

	return resp
}

func (h *Hub) ListenUDP(port int) error {
	udpAddr := net.UDPAddr{}
	udpAddr.Port = port

	conn, connErr := net.ListenUDP("udp", &udpAddr)
	if connErr != nil {
		return connErr
	}

	defer conn.Close()

	msgBuffer := make([]byte, 7)

	for {
		n, addr, err := conn.ReadFromUDP(msgBuffer)
		trimmedMsg := strings.TrimSpace(string(msgBuffer[0:n]))
		fmt.Printf("%s -> %s\n", addr.String(), trimmedMsg)

		if trimmedMsg == "STOP" {
			fmt.Println("Exiting UDP server!")
			return nil
		}

		if err != nil {
			util.HandleErr(err)
			continue
		}

		resp := h.HandleMsgUDP(trimmedMsg, addr)
		_, err = conn.WriteToUDP([]byte(resp), addr) // TODO(antonskwr): handle the number of bytes

		if err != nil {
			util.HandleErr(err)
			continue
		}
	}
}
