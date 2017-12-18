package api

import (
	"../space"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

func getProjectByQuery(query url.Values) (*space.ProjectsEntry, error) {
	qid, ok := query["id"]

	if !ok {
		return nil, fmt.Errorf("Extract ID")
	}

	id, err := strconv.ParseUint(qid[0], 10, 64)

	if err != nil {
		return nil, fmt.Errorf("Invalid ID: %#v", qid)
	}

	project := space.Projects.Get(id)

	if project == nil {
		return nil, fmt.Errorf("Invalid ID: %d", id)
	}

	key := query["key"]
	token, tokenExists := query["token"]

	if tokenExists {
		if !space.Tokens.Has(token[0]) {
			return nil, fmt.Errorf("Invalid TOKEN: %s", token[0])
		}
	} else if project.Key != key[0] {
		return nil, fmt.Errorf("Invalid KEY")
	}

	return project, nil
}

func ProjectGet(w http.ResponseWriter, r *http.Request) {
	project, err := getProjectByQuery(r.URL.Query())

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	end(w, project)
}

func ProjectReg(w http.ResponseWriter, r *http.Request) {
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

	end(w, project)
}
