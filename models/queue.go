package models

import (
	"resque-inspector/resque"
	"unicode"
	"unicode/utf8"
)

const FailedQueue = "failed"

type Queue struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	JobCount int    `json:"job_count"`
	Jobs     []Job  `json:"jobs"`
}

func GetQueueList(filter resque.Filter) resque.Result[Queue] {
	queues := resque.GetList("queues")
	var data []Queue
	filtered := 0
	for _, queue := range queues {
		if resque.ShouldFilterString(filter, queue) {
			filtered++
			continue
		}
		structure := Queue{
			Id:       queue,
			Name:     queue,
			JobCount: int(resque.GetEntryCount(queuePathForRequest(queue))),
			Jobs:     []Job{},
		}

		data = append(data, structure)
	}

	r, size := utf8.DecodeRuneInString(FailedQueue)
	if r == utf8.RuneError { /* no errors possible */
	}
	failedName := string(unicode.ToUpper(r)) + FailedQueue[size:]
	data = append(data, Queue{
		Id:       FailedQueue,
		Name:     failedName,
		JobCount: int(resque.GetEntryCount(queuePathForRequest(FailedQueue))),
		Jobs:     []Job{},
	})

	return resque.Result[Queue]{
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
		JobCount: int(resque.GetEntryCount(queuePathForRequest(name))),
		Jobs:     []Job{},
	}
}

func (q Queue) Clear() error {
	return resque.Clear(q.queuePathForRequest())
}

func (q Queue) IsFailed() bool {
	return q.Id == FailedQueue
}

func (q Queue) queuePathForRequest() string {
	return queuePathForRequest(q.Id)
}

func queuePathForRequest(queue string) string {
	if queue == FailedQueue {
		return queue
	}
	return "queue:" + queue
}
