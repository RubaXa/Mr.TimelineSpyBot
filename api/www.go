package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type ClientEnv struct {
	TS    int64  `json:"ts"`
	Token string `json:"token"`
}

func WwwDashboard(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadFile("./www/index.html")

	envJson, _ := json.Marshal(ClientEnv{
		TS:    time.Now().Unix() * 1000,
		Token: r.URL.Query()["token"][0],
	})

	html := strings.Replace(
		string(body),
		"{/*env*/}",
		string(envJson),
		1,
	)

	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(html))
}
