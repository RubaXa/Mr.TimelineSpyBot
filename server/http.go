package server

import (
	"fmt"
	"net/http"

	"encoding/json"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"net"
	"strings"
	"time"

	"../space"
)

type HttpServer struct {
	Host   string
	events chan wsEvent
	routes Routes
}

type wsEvent struct {
	Id     uint64      `json:"id"`
	Type   string      `json:"type"`
	Detail interface{} `json:"detail"`
}

type wsClient struct {
	project *space.ProjectsEntry
	conn    net.Conn
}

type RouteHandle func(http.ResponseWriter, *http.Request)
type Routes map[string]RouteHandle

func (hs *HttpServer) PushEvent(id uint64, t string, d interface{}) {
	hs.events <- wsEvent{id, t, d}
}

func (hs *HttpServer) Start() error {
	fmt.Println("Start http server", hs.Host)

	for path, handle := range hs.routes {
		http.HandleFunc(path, handle)
	}

	wsClients := make([]wsClient, 0, 1)

	go func() {
		for evt := range hs.events {
			msg, err := json.Marshal(evt)

			if err != nil {
				fmt.Println("[ws] FAILED SEND", err)
				continue
			}

			for _, client := range wsClients {
				if evt.Id == 0 || evt.Id == client.project.Id {
					wsutil.WriteServerMessage(client.conn, ws.OpText, msg)
				}
			}
		}
	}()

	http.HandleFunc("/ws/", func(w http.ResponseWriter, r *http.Request) {
		conn, _, _, err := ws.UpgradeHTTP(r, w, nil)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		parts := strings.Split(r.URL.Path, "/ws/")
		project, err := space.Projects.GetByTokenValue(parts[len(parts)-1])

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		wsClients = append(wsClients, wsClient{project, conn})

		// Ping
		go func() {
			defer func() {
				conn.Close()

				for i, client := range wsClients {
					if client.conn == conn {
						wsClients = append(wsClients[:i], wsClients[i+1:]...)
						break
					}
				}
			}()

			fmt.Println("[ws] START")

			for {
				err := wsutil.WriteServerMessage(conn, ws.OpPing, nil)

				if err != nil {
					fmt.Println("[ws] Ping error:", err)
					break
				}

				time.Sleep(time.Second)
			}

			fmt.Println("[ws] EXIT")
		}()
	})

	return http.ListenAndServe(hs.Host, nil)
}

func CreateHttpServer(h string, r Routes) *HttpServer {
	srv := &HttpServer{
		Host:   h,
		events: make(chan wsEvent, 10),
		routes: r,
	}
	return srv
}
