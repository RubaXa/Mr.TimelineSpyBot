package api

import (
	"../space"
	"net/http"
)

func TokenCreate(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	project, err := getProjectByQuery(query)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	token, err := space.Tokens.Create(project.Id)

	if err != nil {
		http.Error(w, "Failed token create", http.StatusInternalServerError)
		return
	}

	end(w, token)
}
