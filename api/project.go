package api

import (
	"../space"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

func getProjectByQuery(query url.Values) (*space.ProjectsEntry, error) {
	var id uint64

	tokenValue, tokenExists := query["token"]

	if tokenExists {
		token, err := space.Tokens.GetByValue(tokenValue[0])

		if err != nil {
			return nil, err
		}

		id = token.ProjectId
	} else {
		qid, ok := query["id"]

		if !ok {
			return nil, fmt.Errorf("Extract ID")
		}

		pid, err := strconv.ParseUint(qid[0], 10, 64)

		if err != nil {
			return nil, fmt.Errorf("Invalid ID: %#v", qid)
		}

		id = pid
	}

	project := space.Projects.Get(id)

	if project == nil {
		return nil, fmt.Errorf("Invalid ID: %d", id)
	}

	if !tokenExists {
		key := query["key"]

		if project.Key != key[0] {
			return nil, fmt.Errorf("Invalid KEY")
		}
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
