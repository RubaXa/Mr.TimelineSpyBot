package server

import (
	"fmt"
	"net/http"
)

type HttpServer struct {
	Host string
}

type RouteHandle func(http.ResponseWriter, *http.Request)
type Routes map[string]RouteHandle

func (hs *HttpServer) Start(routes Routes) error {
	fmt.Println("Start http server", hs.Host)

	for path, handle := range routes {
		http.HandleFunc(path, handle)
	}

	return http.ListenAndServe(hs.Host, nil)
}
