package models

import (
	"encoding/json"
	"fmt"
	"resque-inspector/resque"
	"time"
)

type JobInterface interface {
	Stringify() string
}
type Job struct {
	Class     string                   `json:"class"`
	Args      []map[string]interface{} `json:"args"`
	Id        string                   `json:"id"`
	Prefix    string                   `json:"prefix"`
	QueueTime float64                  `json:"queue_time"`
}

func (f Job) Stringify() string {
	return fmt.Sprintf("class: %s", f.Class)
}

type FailedJob struct {
	Payload   Job       `json:"payload"`
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

func (q Queue) GetJobList(filter resque.Filter, start int64, limit int64) resque.Result[JobInterface] {
	var entries []string
	var classes = make([]string, 0)
	var exceptions = make([]string, 0)
	var data = make([]JobInterface, 0)

	if q.Id == "failed" {
		entries = resque.GetEntries(q.Id, start, limit)
	} else {
		entries = resque.GetEntries("queue:"+q.Id, start, limit)
	}

	for _, entry := range entries {
		if q.Id == "failed" {
			var job FailedJob
			err := json.Unmarshal([]byte(entry), &job)
			if err != nil {
				continue
			}
			if ShouldFilterFailedJob(filter, job) {
				continue
			}

			classes = append(classes, job.Payload.Class)
			exceptions = append(exceptions, job.Exception)
			data = append(data, job)
			continue
		}

		var job Job
		err := json.Unmarshal([]byte(entry), &job)
		if err != nil {
			continue
		}
		if ShouldFilterJob(filter, job) {
			continue
		}

		classes = append(classes, job.Class)
		data = append(data, job)
	}

	return resque.Result[JobInterface]{
		Filter:     filter,
		Filtered:   filter.Filtered,
		Total:      len(data),
		Classes:    classes,
		Exceptions: exceptions,
		Items:      data,
	}
}

func ShouldFilterJob(f resque.Filter, job Job) bool {
	if f.Class == "" {
		return false
	}
	if f.Class == job.Class {
		return true
	}
	return false
}

func ShouldFilterFailedJob(f resque.Filter, failed FailedJob) bool {
	if f.Class == "" && f.Exception == "" {
		return false
	}
	if f.Class == failed.Payload.Class {
		return true
	}
	if f.Exception == failed.Exception {
		return true
	}
	return false
}

func (f FailedJob) Retry() error {
	out, err := json.Marshal(f.Payload)
	if err != nil {
		return err
	}
	return resque.Retry(f.Queue, string(out))
}
