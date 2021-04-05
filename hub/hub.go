package hub

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"sync"

	"github.com/antonskwr/nat-punch-through-hub/util"
)

type RespType int

const (
	RespTypeNone RespType = iota
	RespTypeRegular
	RespTypePeerReq
)

type Resp struct {
	rType      RespType
	msg        string
	peerAddr   *net.UDPAddr
	peerReqMsg string
}

func ResponseNone() Resp {
	return Resp{
		RespTypeNone,
		"",
		nil,
		"",
	}
}

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

func (h *Hub) HandleMsgUDP(msg string, addr *net.UDPAddr) Resp {
	resp := ResponseNone()
	splittedMsgs := strings.Split(msg, " ")

	switch true {
	case len(splittedMsgs) == 1:
		switch splittedMsgs[0] {
		case "HB":
			return resp
		case "LIST":
			resp.rType = RespTypeRegular
			if len(h.gameServers) == 0 {
				resp.msg = "No hosts registered"
				return resp // TODO(antonskwr): dirty code
			}
			resp.msg = "Listing hosts ...\n"
			for id := range h.gameServers {
				resp.msg = fmt.Sprintf("%sHost[%d]\n", resp.msg, id)
			}
		}
	case len(splittedMsgs) == 2:
		id, err := strconv.Atoi(splittedMsgs[1])
		if err == nil {
			switch splittedMsgs[0] {
			case "ADD":
				h.listLock.Lock()
				h.gameServers[id] = GameServer{addr}
				h.listLock.Unlock()
				resp.rType = RespTypeRegular
				resp.msg = "Host added successfully"
			case "JOIN":
				gameServer, ok := h.gameServers[id]
				if !ok {
					resp.rType = RespTypeRegular
					resp.msg = "No matching ids"
				} else {
					if gameServer.addr.String() != addr.String() {
						resp.rType = RespTypePeerReq
						resp.msg = "OK " + gameServer.addr.String()
						resp.peerAddr = gameServer.addr
						resp.peerReqMsg = "REQ " + addr.String()
					} else {
						resp.rType = RespTypeRegular
						resp.msg = "can't connect"
					}
				}
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

		if addr == nil {
			err := fmt.Errorf("ListenUDP: ReadFromUDP addr is nil")
			util.HandleErrNonFatal(err)
			continue
		}

		trimmedMsg := strings.TrimSpace(string(msgBuffer[0:n]))

		if err != nil {
			util.HandleErrNonFatal(err)
			continue
		}

		resp := h.HandleMsgUDP(trimmedMsg, addr)

		if resp.rType != RespTypeNone {
			fmt.Printf("%s -> %s\n", addr.String(), trimmedMsg)
			_, err = conn.WriteToUDP([]byte(resp.msg), addr) // TODO(antonskwr): handle the number of bytes
			if err != nil {
				util.HandleErrNonFatal(err)
				continue
			}

			switch resp.rType {
			case RespTypeRegular:
			case RespTypePeerReq:
				pAddr := resp.peerAddr
				pReqMsg := resp.peerReqMsg
				if pAddr != nil && len(pReqMsg) > 4 {
					_, err = conn.WriteToUDP([]byte(pReqMsg), pAddr) // TODO(antonskwr): handle the number of bytes
					fmt.Println("Sent a request to host")
					if err != nil {
						util.HandleErrNonFatal(err)
						continue
					}
				} else {
					fmt.Printf("request addr is nil:%v, req len %v\n", pAddr == nil, len(pReqMsg))
				}
			}
		}
	}
}
