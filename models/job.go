package models

import (
	"encoding/json"
	"fmt"
	"resque-inspector/resque"
	"resque-inspector/result"
	"time"
)

type JobInterface interface {
	Stringify() string
}
type Job struct {
	Payload struct {
		Class     string                   `json:"class"`
		Args      []map[string]interface{} `json:"args"`
		Id        string                   `json:"id"`
		Prefix    string                   `json:"prefix"`
		QueueTime float64                  `json:"queue_time"`
	} `json:"payload"`
}

func (f Job) Stringify() string {
	return fmt.Sprintf("class: %s", f.Payload.Class)
}

type FailedJob struct {
	Job
	FailedAt  time.Time `json:"failed_at"`
	Exception string    `json:"exception"`
	Error     string    `json:"error"`
	Backtrace []string  `json:"backtrace"`
	Worker    string    `json:"worker"`
	Queue     string    `json:"queue"`
}

func (f FailedJob) Stringify() string {
	return fmt.Sprintf("error: %s\n\texception: %s\n\tqueue: %s\n", f.Error, f.Exception, f.Queue)
}

func (q Queue) GetJobList(filter string) result.Result[JobInterface] {
	resque.PrepareClient()
	defer resque.Client.Close()

	var filtered = 0
	var entries []string
	var data []JobInterface

	if q.Id == "failed" {
		entries = resque.GetEntries(q.Id, true)
	} else {
		entries = resque.GetEntries("queue:"+q.Id, true)
	}

	for _, entry := range entries {
		if q.Id == "failed" {
			var job FailedJob
			err := json.Unmarshal([]byte(entry), &job)
			if err != nil {
				filtered++
				continue
			}

			data = append(data, job)
			continue
		}

		var job Job
		err := json.Unmarshal([]byte(entry), &job)
		if err != nil {
			filtered++
			continue
		}

		data = append(data, job)
	}

	return result.Result[JobInterface]{
		Filter:   filter,
		Filtered: filtered,
		Total:    len(data),
		Items:    data,
	}
}
