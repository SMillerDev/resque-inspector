package models

import (
	"encoding/json"
	"regexp"
	"resque-inspector/resque"
	"resque-inspector/result"
	"strings"
)

type WorkerSlot struct {
	id     string
	Name   string `json:"name"`
	Host   string `json:"host"`
	Socket string `json:"socket"`
	Entry  Job    `json:"entry"`
}

func GetWorkerList(filter string) result.NamedResult[WorkerSlot] {
	workers := resque.GetList("workers", true)
	data := make(map[string][]WorkerSlot)
	var filtered = 0
	for worker := range workers {
		match, _ := regexp.MatchString(filter, worker)
		if !match {
			filtered++
			continue
		}

		parts := strings.Split(worker, ":")
		if len(data[parts[2]]) == 0 {
			data[parts[2]] = []WorkerSlot{}
		}
		structure := WorkerSlot{
			id:     worker,
			Name:   parts[2],
			Host:   parts[0],
			Socket: parts[1],
		}
		var job Job
		err := json.Unmarshal([]byte(resque.GetEntry("worker:"+worker, true)), &job)
		if err != nil {
			continue
		}
		structure.Entry = job

		data[parts[2]] = append(data[parts[2]], structure)
	}

	return result.NamedResult[WorkerSlot]{
		Filter:   filter,
		Total:    len(data),
		Filtered: filtered,
		Items:    data,
	}
}
