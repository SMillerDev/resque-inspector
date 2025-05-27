package models

import (
	"encoding/json"
	"fmt"
	"log"
	"resque-inspector/resque"
	"time"
)

type JobInterface interface {
	Stringify() string
	PayloadString() string
	Identifier() string
	QueueIdentifier() string
}
type Job struct {
	Class     string                   `json:"class"`
	Args      []map[string]interface{} `json:"args"`
	Id        string                   `json:"id"`
	Prefix    string                   `json:"prefix"`
	QueueTime float64                  `json:"queue_time"`
	RetryTime float64                  `json:"retry_time"`
}

func (f Job) Stringify() string {
	return fmt.Sprintf("class: %s", f.Class)
}
func (f Job) Identifier() string      { return f.Id }
func (f Job) QueueIdentifier() string { return "" }
func (f Job) PayloadString() string {
	f.RetryTime = float64(time.Now().Unix())
	str, _ := json.Marshal(f)

	return string(str)
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
	return fmt.Sprintf("error: %s\n\texception: %s\n\tqueue: %s\n%s", f.Error, f.Exception, f.Queue, f.Payload.Stringify())
}
func (f FailedJob) QueueIdentifier() string { return f.Queue }
func (f FailedJob) Identifier() string      { return f.Payload.Identifier() }
func (f FailedJob) PayloadString() string   { return f.Payload.PayloadString() }

func (q Queue) GetJobList(filter resque.Filter, start int64, limit int64) resque.Result[JobInterface] {
	var entries []string
	var classes = make(map[string]int)
	var exceptions = make(map[string]int)
	var data = make([]JobInterface, 0)
	filtered := 0

	if q.JobCount == 0 && q.Id != "failed" {
		return resque.Result[JobInterface]{
			Filter:     filter,
			Filtered:   filtered,
			Total:      q.JobCount,
			Selected:   0,
			Classes:    classes,
			Exceptions: exceptions,
			Items:      data,
		}
	}

	if q.Id == "failed" {
		entries = resque.GetEntries(q.Id, start, limit)
	} else {
		entries = resque.GetEntries("queue:"+q.Id, start, limit)
	}

	if len(entries) == 0 {
		return resque.Result[JobInterface]{
			Filter:     filter,
			Filtered:   filtered,
			Total:      q.JobCount,
			Selected:   0,
			Classes:    classes,
			Exceptions: exceptions,
			Items:      data,
		}
	}

	for _, entry := range entries {
		if q.Id == "failed" {
			var job FailedJob
			err := json.Unmarshal([]byte(entry), &job)
			if err != nil {
				continue
			}
			if ShouldFilterFailedJob(filter, job) {
				filtered++
				continue
			}

			classes[job.Payload.Class]++
			exceptions[job.Exception]++
			data = append(data, job)
			continue
		}

		var job Job
		err := json.Unmarshal([]byte(entry), &job)
		if err != nil {
			continue
		}
		if ShouldFilterJob(filter, job) {
			filtered++
			continue
		}

		classes[job.Class]++
		data = append(data, job)
	}

	return resque.Result[JobInterface]{
		Filter:     filter,
		Filtered:   filtered,
		Total:      q.JobCount,
		Selected:   len(data),
		Classes:    classes,
		Exceptions: exceptions,
		Items:      data,
	}
}

func ShouldFilterJob(f resque.Filter, job Job) bool {
	if f.Class != "" && f.Class != job.Class {
		if Debug {
			log.Default().Println("Filter job class does not match.")
		}

		return true
	}

	if f.Id != "" && f.Id != job.Id {
		if Debug {
			log.Default().Printf("[Filter] job id (%s) does not match (%s).)", job.Id, f.Id)
		}
		return true
	}

	if Debug {
		log.Default().Println("Filter matched.")
	}
	return false
}

func ShouldFilterFailedJob(f resque.Filter, failed FailedJob) bool {
	if f.Class != "" && f.Class != failed.Payload.Class {
		if Debug {
			log.Default().Printf("[Filter] job class (%s) does not match (%s).)", failed.Payload.Class, f.Class)
		}
		return true
	}
	if f.Exception != "" && f.Exception != failed.Exception {
		if Debug {
			log.Default().Printf("[Filter] job exception (%s) does not match (%s).)", failed.Exception, f.Exception)
		}
		return true
	}

	if f.Id != "" && f.Id != failed.Payload.Id {
		if Debug {
			log.Default().Printf("[Filter] job id (%s) does not match (%s).)", failed.Payload.Id, f.Id)
		}
		return true
	}

	if Debug {
		log.Default().Println("Filter matched.")
	}

	return false
}

func (f FailedJob) Retry() error {
	out, err := json.Marshal(f.Payload)
	if err != nil {
		return err
	}
	return resque.Queue(f.Queue, string(out))
}
