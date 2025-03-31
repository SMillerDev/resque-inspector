package models

import (
	"resque-inspector/resque"
	"resque-inspector/result"
)

type Queue struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	JobCount int64  `json:"job_count"`
	Jobs     []Job  `json:"jobs"`
}

func GetQueueList(filter result.Filter) result.Result[Queue] {
	queues := resque.GetList("queues", true)
	var data []Queue
	filtered := 0
	for _, queue := range queues {
		if result.ShouldFilterString(filter, queue) {
			continue
		}
		structure := Queue{
			Id:       queue,
			Name:     queue,
			JobCount: resque.GetEntryCount("queue:"+queue, true),
			Jobs:     []Job{},
		}

		data = append(data, structure)
	}
	data = append(data, Queue{
		Id:       "failed",
		Name:     "Failed",
		JobCount: resque.GetEntryCount("failed", true),
		Jobs:     []Job{},
	})

	return result.Result[Queue]{
		Filter:   filter,
		Total:    len(data),
		Filtered: filtered,
		Items:    data,
	}
}

func GetQueue(name string) Queue {
	return Queue{
		Id:       name,
		Name:     name,
		JobCount: resque.GetEntryCount(name, true),
		Jobs:     []Job{},
	}
}
