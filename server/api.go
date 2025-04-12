package server

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"resque-inspector/models"
	"strconv"
)

func getRootApi(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("got %s API request\n", r.RequestURI)
	typeVal := r.PathValue("type")

	var jsonData []byte
	switch typeVal {
	case "queues":
		result := models.GetQueueList(filterFromRequest(r))
		out, err := json.Marshal(result)
		if err != nil {
			log.Default().Printf("could not marshal json: %s\n", err)
			returnError(w, http.StatusInternalServerError, map[string]interface{}{})
			return
		}
		jsonData = out
	case "workers":
		result := models.GetWorkerList(filterFromRequest(r))
		out, err := json.Marshal(result)
		if err != nil {
			log.Default().Printf("could not marshal json: %s\n", err)
			returnError(w, http.StatusInternalServerError, map[string]interface{}{})
			return
		}
		jsonData = out
	default:
		log.Default().Printf("received unknown API request: %s\n", typeVal)
		returnError(w, http.StatusBadRequest, map[string]interface{}{})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, _ = io.WriteString(w, string(jsonData))
}

func getJobsApi(w http.ResponseWriter, r *http.Request) {
	queueVal := r.PathValue("queue")
	if queueVal == "" {
		log.Default().Printf("received unknown API request: %s\n", r.RequestURI)
		returnError(w, http.StatusBadRequest, map[string]interface{}{})
		return
	}
	start := 0
	if r.URL.Query().Has("start") {
		start, _ = strconv.Atoi(r.URL.Query().Get("start"))
	}
	end := 100
	if r.URL.Query().Has("offset") {
		end, _ = strconv.Atoi(r.URL.Query().Get("offset"))
	}

	result := models.GetQueue(queueVal).GetJobList(filterFromRequest(r), int64(start), int64(end))
	out, err := json.Marshal(result)
	if err != nil {
		log.Default().Println(result)
		log.Default().Printf("could not marshal json: %s\n", err)
		returnError(w, http.StatusInternalServerError, map[string]interface{}{})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, _ = io.WriteString(w, string(out))
}

func retryJobApi(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		returnError(w, http.StatusMethodNotAllowed, map[string]interface{}{})
		return
	}

	queueVal := r.PathValue("queue")
	idVal := r.PathValue("id")
	if queueVal == "" || idVal == "" {
		log.Default().Printf("received unknown API request: %s\n", r.RequestURI)
		returnError(w, http.StatusBadRequest, map[string]interface{}{})
		return
	}
}

func clearApi(w http.ResponseWriter, r *http.Request) {
	if r.Method != "DELETE" {
		returnError(w, http.StatusMethodNotAllowed, map[string]interface{}{})
		return
	}
	queueVal := r.PathValue("queue")
	if queueVal == "" {
		log.Default().Printf("received unknown API request: %s\n", r.RequestURI)
		returnError(w, http.StatusBadRequest, map[string]interface{}{})
		return
	}

	err := models.GetQueue(queueVal).Clear()
	if err != nil {
		log.Default().Printf("could not clear queue: %s\n", err)
		returnError(w, http.StatusInternalServerError, map[string]interface{}{})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, _ = io.WriteString(w, string("{}"))
}
