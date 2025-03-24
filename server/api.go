package server

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"resque-inspector/models"
)

func GetApi(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("got %s request\n", r.RequestURI)
	typeVal := r.PathValue("type")
	fmt.Println("Request type: ", typeVal)

	var jsonData []byte
	switch typeVal {
	case "queues":
		result := models.GetQueueList(getFilter(r))
		out, err := json.Marshal(result)
		if err != nil {
			log.Default().Printf("could not marshal json: %s\n", err)
			returnError(w, http.StatusInternalServerError, map[string]interface{}{})
			return
		}
		jsonData = out
	case "workers":
		result := models.GetWorkerList(getFilter(r))
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
