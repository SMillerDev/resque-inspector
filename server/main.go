package server

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"resque-inspector/resque"
	"time"
)

const httpAddr = ":5678"

func returnError(w http.ResponseWriter, code int, data interface{}) {
	w.WriteHeader(code)
	w.Header().Set("Content-Type", "application/json")
	jsonData, _ := json.Marshal(data)
	_, _ = io.WriteString(w, string(jsonData))
}

func filterFromRequest(r *http.Request) resque.Filter {
	return resque.Filter{
		Regex:     r.URL.Query().Get("filter"),
		Class:     "",
		Exception: "",
		Queue:     "",
		StartDate: time.Time{},
		EndDate:   time.Time{},
		Filtered:  0,
	}
}

func Serve() {
	http.HandleFunc("/{page}", getUi)
	http.HandleFunc("/", getUi)

	http.HandleFunc("/api/v1/{type}", getRootApi)
	http.HandleFunc("/api/v1/queues/{queue}/jobs", getJobsApi)

	err := http.ListenAndServe(httpAddr, nil)
	if errors.Is(err, http.ErrServerClosed) {
		log.Default().Printf("server closed\n")
	} else if err != nil {
		log.Default().Fatalf("error starting server: %s\n", err)
	}
}
