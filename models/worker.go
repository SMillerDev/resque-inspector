package models

import (
	"encoding/json"
	"resque-inspector/resque"
	"strings"
)

type WorkerSlot struct {
	id     string
	Name   string `json:"name"`
	Host   string `json:"host"`
	Socket string `json:"socket"`
	Entry  Job    `json:"entry"`
}

func GetWorkerList(filter resque.Filter) resque.NamedResult[WorkerSlot] {
	workers := resque.GetList("workers")
	data := make(map[string][]WorkerSlot)
	var filtered = 0
	for _, worker := range workers {
		if resque.ShouldFilterString(filter, worker) {
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
		entry := resque.GetEntryOrNil("worker:" + worker)
		if entry != "" {
			err := json.Unmarshal([]byte(entry), &job)
			if err != nil {
				continue
			}
			structure.Entry = job
		}

		data[parts[2]] = append(data[parts[2]], structure)
	}

	return resque.NamedResult[WorkerSlot]{
		Filter:   filter,
		Total:    len(data),
		Filtered: filtered,
		Items:    data,
	}
}
