package hub

import (
	"fmt"
	"log"
	"net/http"
)

type Hub struct {
	mux *http.ServeMux
}

func NewHub() *Hub {
	h := &Hub{
		mux: http.NewServeMux(),
	}

	h.setupRoutes()
	return h
}

func (h *Hub) ListenAndServe(port string) error {
	// TODO(antonskwr): change for secured TLS connection (ListenAndServeTLS)
	return http.ListenAndServe(port, h.mux)
}

func (h *Hub) setupRoutes() {
	h.Handle("/", http.NotFoundHandler())
	h.Handle("/favicon.ico", http.NotFoundHandler())

	h.Handle("/api/servers/list", h.handleListServers())
	h.Handle("/api/servers/add", h.handleAddServer())
	h.Handle("/api/servers/connect", h.handleConnectPeerToServer())
}

func (h *Hub) Handle(pattern string, handler http.Handler) {
	h.mux.Handle(pattern, handler)
}

func handleErr(err error, message ...string) {
	if err != nil {
		if len(message) > 0 {
			err = fmt.Errorf("[%s] -- %w --", message[0], err)
		}
		log.Fatal(err)
	}
}
