package api

import (
	"../space"
	"encoding/json"
	"net/http"
	"strconv"
)

func Get(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(r.URL.Query()["id"][0], 10, 64)

	if err != nil {
		http.Error(w, "Invalid ID", http.StatusInternalServerError)
	}

	project := space.Projects.Get(id)

	if project.Key != r.URL.Query()["key"][0] {
		http.Error(w, "Invalid KEY", http.StatusForbidden)
		return
	}

	jsonBody, err := json.Marshal(project)

	if err != nil {
		http.Error(w, "Error converting results to json", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonBody)
}

func Reg(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()

	if err != nil {
		http.Error(w, "Error request.ParseForm", http.StatusInternalServerError)
		return
	}

	q := r.URL.Query()
	project := space.Projects.New()
	project.Name = q["name"][0]

	err = space.Projects.Save(project)

	if err != nil {
		http.Error(w, "Failed save", http.StatusInternalServerError)
		return
	}

	jsonBody, err := json.Marshal(project)

	if err != nil {
		http.Error(w, "Error converting results to json", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonBody)
}
