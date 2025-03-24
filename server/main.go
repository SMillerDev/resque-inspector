package server

import (
	"encoding/json"
	"io"
	"net/http"
)

var Filter string

func returnError(w http.ResponseWriter, code int, data interface{}) {
	w.WriteHeader(code)
	w.Header().Set("Content-Type", "application/json")
	jsonData, _ := json.Marshal(data)
	_, _ = io.WriteString(w, string(jsonData))
}

func getFilter(r *http.Request) string {
	r.URL.Query().Get("filter")
	if Filter == "" {
		Filter = ".*"
	}

	return Filter
}
