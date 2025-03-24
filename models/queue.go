package models

import (
	"regexp"
	"resque-inspector/resque"
	"resque-inspector/result"
)

type Queue struct {
	Id       string
	Name     string `json:"name"`
	JobCount int64  `json:"job_count"`
	Jobs     []Job  `json:"jobs"`
}

func GetQueueList(filter string) result.Result[Queue] {
	queues := resque.GetList("queues", true)
	var data []Queue
	filtered := 0
	for queue := range queues {
		match, _ := regexp.MatchString(filter, queue)
		if !match {
			filtered++
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
