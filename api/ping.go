package api

import (
	"net/http"
	"time"
)

func Ping(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte("{\"pong\":\"" + time.Now().String() + "\"}"))
}
