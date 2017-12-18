package api

import (
	"../space"
	"math"
	"net/http"
)

func RecordList(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	project, err := getProjectByQuery(q)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	records, err := space.Records.Select(project.Id, 0, math.MaxInt32)

	if err != nil {
		http.Error(w, "Error Records.Select", http.StatusInternalServerError)
		return
	}

	end(w, records)
}
