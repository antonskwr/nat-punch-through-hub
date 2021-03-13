package hub

import (
	"fmt"
	"net"
	"net/http"
)

func printIpAndPort(name, addr string) {
	fmt.Printf("%v: ", name)
	ip, port, err := net.SplitHostPort(addr)
	if err != nil {
		fmt.Println()
		return
	}

	fmt.Printf("%v %v\n", ip, port)
}

func printSeparator() {
	fmt.Printf("=========\n\n")
}

func (h *Hub) handleConnectPeerToServer() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			IPAddressReal := r.Header.Get("X-Real-Ip")
			IPAddressForwared := r.Header.Get("X-Forwarded-For")
			IPAddressRemote := r.RemoteAddr

			printIpAndPort("IPAddressReal", IPAddressReal)
			printIpAndPort("IPAddressForwared", IPAddressForwared)
			printIpAndPort("IPAddressRemote", IPAddressRemote)
			printSeparator()

			// if IPAddress == "" {
			// 		IPAddress = r.Header.Get("X-Forwarded-For")
			// }
			// if IPAddress == "" {
			// 		IPAddress = r.RemoteAddr
			// }

			w.Write([]byte("handleConnectPeerToServer"))
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}
}

func (h *Hub) handleListServers() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:

			IPAddressReal := r.Header.Get("X-Real-Ip")
			IPAddressForwared := r.Header.Get("X-Forwarded-For")
			IPAddressRemote := r.RemoteAddr

			printIpAndPort("IPAddressReal", IPAddressReal)
			printIpAndPort("IPAddressForwared", IPAddressForwared)
			printIpAndPort("IPAddressRemote", IPAddressRemote)
			printSeparator()

			// if IPAddress == "" {
			// 		IPAddress = r.Header.Get("X-Forwarded-For")
			// }
			// if IPAddress == "" {
			// 		IPAddress = r.RemoteAddr
			// }

			w.Write([]byte("handleListServers"))
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}
}

func (h *Hub) handleAddServer() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}
}
