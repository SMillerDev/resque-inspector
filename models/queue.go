package models

import (
	"resque-inspector/resque"
)

type Queue struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	JobCount int    `json:"job_count"`
	Jobs     []Job  `json:"jobs"`
}

func GetQueueList(filter resque.Filter) resque.Result[Queue] {
	queues := resque.GetList("queues")
	var data []Queue
	for _, queue := range queues {
		if resque.ShouldFilterString(filter, queue) {
			continue
		}
		structure := Queue{
			Id:       queue,
			Name:     queue,
			JobCount: int(resque.GetEntryCount("queue:" + queue)),
			Jobs:     []Job{},
		}

		data = append(data, structure)
	}
	data = append(data, Queue{
		Id:       "failed",
		Name:     "Failed",
		JobCount: int(resque.GetEntryCount("failed")),
		Jobs:     []Job{},
	})

	return resque.Result[Queue]{
		Filter:   filter,
		Total:    len(data),
		Filtered: filter.Filtered,
		Items:    data,
	}
}

func GetQueue(name string) Queue {
	return Queue{
		Id:       name,
		Name:     name,
		JobCount: int(resque.GetEntryCount("queue:" + name)),
		Jobs:     []Job{},
	}
}

func (q Queue) Clear() error {
	return resque.Clear("queue:" + q.Id)
}
