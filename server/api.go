package server

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"resque-inspector/models"
	"resque-inspector/resque"
	"strconv"
)

func getRootApi(w http.ResponseWriter, r *http.Request) {
	log.Default().Printf("[API] %s %s request\n", r.Method, r.RequestURI)
	if r.Method != "GET" {
		returnError(w, http.StatusMethodNotAllowed, map[string]interface{}{})
		return
	}
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
	log.Default().Printf("[API] %s %s request\n", r.Method, r.RequestURI)
	if r.Method != "GET" {
		returnError(w, http.StatusMethodNotAllowed, map[string]interface{}{})
		return
	}
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

func modifyJobApi(w http.ResponseWriter, r *http.Request) {
	log.Default().Printf("[API] %s %s request\n", r.Method, r.RequestURI)
	if r.Method != "POST" && r.Method != "DELETE" {
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
	result := models.GetQueue(queueVal).GetJobList(resque.Filter{Id: idVal}, 0, 1000)
	if len(result.Items) < 1 {
		returnError(w, http.StatusNotFound, map[string]interface{}{})
		return
	}

	if r.Method == "POST" {
		err := resque.Queue(result.Items[0].QueueIdentifier(), result.Items[0].PayloadString())
		if err != nil {
			returnError(w, http.StatusInternalServerError, map[string]interface{}{})
		}
	}
	if r.Method == "DELETE" {
		err := resque.Delete(queueVal, result.Items[0].Identifier())
		if err != nil {
			returnError(w, http.StatusInternalServerError, map[string]interface{}{})
		}
	}
}

func clearApi(w http.ResponseWriter, r *http.Request) {
	log.Default().Printf("[API] %s %s request\n", r.Method, r.RequestURI)
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

	w.WriteHeader(http.StatusNoContent)
}
