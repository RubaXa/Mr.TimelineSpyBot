package api

import (
	"encoding/json"
	"net/http"
)

func init() {
}

func end(w http.ResponseWriter, v interface{}) {
	jsonBody, err := json.Marshal(v)

	if err != nil {
		http.Error(w, "Error converting results to json: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonBody)
}
