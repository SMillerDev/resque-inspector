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
	Id        string
	Regex     string
	Class     string
	Exception string
	Queue     string
	StartDate time.Time
	EndDate   time.Time
}

func ShouldFilterString(f Filter, queue string) bool {
	if f.Regex == "" {
		return false
	}
	matches, _ := regexp.MatchString(f.Regex, queue)
	return !matches
}
