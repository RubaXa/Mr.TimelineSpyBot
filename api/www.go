package api

import (
	"io/ioutil"
	"net/http"
)

func WwwDashboard(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadFile("./www/index.html")

	w.Header().Set("Content-Type", "text/html")
	w.Write(body)
}
