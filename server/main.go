package server

import (
	"encoding/json"
	"errors"
	"github.com/NYTimes/gziphandler"
	"io"
	"log"
	"net/http"
	"resque-inspector/resque"
	"strconv"
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
	start, _ := strconv.Atoi(r.URL.Query().Get("startDate"))
	end, _ := strconv.Atoi(r.URL.Query().Get("endDate"))
	return resque.Filter{
		Regex:     r.URL.Query().Get("regex"),
		Class:     r.URL.Query().Get("class"),
		Exception: r.URL.Query().Get("exception"),
		Queue:     r.URL.Query().Get("queue"),
		StartDate: time.UnixMilli(int64(start)),
		EndDate:   time.UnixMilli(int64(end)),
	}
}

func Serve() {
	http.Handle("/{page}", gziphandler.GzipHandler(http.HandlerFunc(getUi)))
	http.Handle("/", gziphandler.GzipHandler(http.HandlerFunc(getUi)))

	http.Handle("/api/v1/{type}", gziphandler.GzipHandler(http.HandlerFunc(getRootApi)))
	http.Handle("/api/v1/queues/{queue}/jobs", gziphandler.GzipHandler(http.HandlerFunc(getJobsApi)))
	http.Handle("/api/v1/queues/{queue}/jobs/{id}", gziphandler.GzipHandler(http.HandlerFunc(modifyJobApi)))
	http.Handle("/api/v1/queues/{queue}", gziphandler.GzipHandler(http.HandlerFunc(clearApi)))

	err := http.ListenAndServe(httpAddr, nil)
	if errors.Is(err, http.ErrServerClosed) {
		log.Default().Printf("server closed\n")
	} else if err != nil {
		log.Default().Fatalf("error starting server: %s\n", err)
	}
}
