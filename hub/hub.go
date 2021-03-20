package hub

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"sync"

	"github.com/antonskwr/nat-punch-through-hub/util"
)

type Hub struct {
	listLock    sync.RWMutex
	gameServers map[int]GameServer
}

func NewHub() *Hub {
	h := &Hub{
		gameServers: make(map[int]GameServer),
	}
	return h
}

type GameServer struct {
	addr *net.UDPAddr
}

func (h *Hub) HandleMsgUDP(msg string, addr *net.UDPAddr) (string, *net.UDPAddr, string) {
	resp := "Unknown msg"
	splittedMsgs := strings.Split(msg, " ")

	if len(splittedMsgs) == 1 {
		if splittedMsgs[0] == "LIST" {
			if len(h.gameServers) == 0 {
				return "No hosts registered", nil, ""
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
				h.gameServers[id] = GameServer{addr}
				h.listLock.Unlock()
				resp = "Host added successfully"
			case "JOIN":
				gameServer, ok := h.gameServers[id]
				if !ok {
					resp = "No matching ids"
				} else {
					if gameServer.addr.String() != addr.String() {
						resp = "OK " + gameServer.addr.String()
						return resp, gameServer.addr, "REQ " + addr.String()
					} else {
						resp = "can't connect"
					}
				}
				// TODO(antonskwr): start NAT punch through
			}
		}
	}

	return resp, nil, ""
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

		if addr == nil {
			err := fmt.Errorf("ListenUDP: ReadFromUDP addr is nil")
			util.HandleErrNonFatal(err)
			continue
		}

		trimmedMsg := strings.TrimSpace(string(msgBuffer[0:n]))
		fmt.Printf("%s -> %s\n", addr.String(), trimmedMsg)

		if err != nil {
			util.HandleErrNonFatal(err)
			continue
		}

		resp, rAddr, req := h.HandleMsgUDP(trimmedMsg, addr)
		_, err = conn.WriteToUDP([]byte(resp), addr) // TODO(antonskwr): handle the number of bytes

		if err != nil {
			util.HandleErrNonFatal(err)
			continue
		}

		if rAddr != nil && len(req) > 4 {
			_, err = conn.WriteToUDP([]byte(req), rAddr) // TODO(antonskwr): handle the number of bytes

			if err != nil {
				util.HandleErrNonFatal(err)
				continue
			}
		}
	}
}
