package resque

import (
	"regexp"
	"time"
)

type Result[T any] struct {
	Filter     Filter         `json:"filter"`
	Total      int            `json:"total"`
	Selected   int            `json:"selected"`
	Filtered   int            `json:"filtered"`
	Classes    map[string]int `json:"classes"`
	Exceptions map[string]int `json:"exceptions"`
	Items      []T            `json:"items"`
}

type NamedResult[T any] struct {
	Filter   Filter         `json:"filter"`
	Total    int            `json:"total"`
	Filtered int            `json:"filtered"`
	Items    map[string][]T `json:"items"`
}

type Filter struct {
	Id        string    `json:"id,omitempty"`
	Regex     string    `json:"regex,omitempty"`
	Class     string    `json:"class,omitempty"`
	Exception string    `json:"exception,omitempty"`
	Queue     string    `json:"queue,omitempty"`
	StartDate time.Time `json:"start_date,omitempty"`
	EndDate   time.Time `json:"end_date,omitempty"`
}

func ShouldFilterString(f Filter, queue string) bool {
	if f.Regex == "" {
		return false
	}
	matches, _ := regexp.MatchString(f.Regex, queue)
	return !matches
}
