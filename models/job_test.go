package models

import (
	"regexp"
	"testing"
	"time"
)

// TestJob calls Job.Stringify() checking for a valid format
func TestJob(t *testing.T) {
	job := Job{
		Class:     "Class",
		Args:      make([]map[string]interface{}, 0),
		Id:        "some-id",
		Prefix:    "SomePrefix",
		QueueTime: 1110000,
	}
	want := regexp.MustCompile(`class: .*`)
	msg := job.Stringify()
	if !want.MatchString(msg) {
		t.Errorf(`Job("Test") = %q, want match for %#q, nil`, msg, want)
	}
}

// TestJob calls Job.Stringify() checking for a valid format
func TestFailedJob(t *testing.T) {
	job := Job{
		Class:     "Class",
		Args:      make([]map[string]interface{}, 0),
		Id:        "some-id",
		Prefix:    "SomePrefix",
		QueueTime: 1110000,
	}
	fail := FailedJob{
		Payload:   job,
		FailedAt:  time.Time{},
		Exception: "A",
		Error:     "B",
		Backtrace: make([]string, 0),
		Worker:    "C",
		Queue:     "D",
	}
	want := regexp.MustCompile(`error: B\n\texception: A\n\tqueue: D`)
	msg := fail.Stringify()
	if !want.MatchString(msg) {
		t.Errorf(`FailedJob("Test") = %q, want match for %#q, nil`, msg, want)
	}
}
