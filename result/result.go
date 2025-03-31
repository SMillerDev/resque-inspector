package result

import (
	"regexp"
	"time"
)

type Result[T any] struct {
	Filter     Filter   `json:"filter"`
	Total      int      `json:"total"`
	Filtered   int      `json:"filtered"`
	Classes    []string `json:"classes,omitempty"`
	Exceptions []string `json:"exceptions,omitempty"`
	Items      []T      `json:"items"`
}

type NamedResult[T any] struct {
	Filter   Filter         `json:"filter"`
	Total    int            `json:"total"`
	Filtered int            `json:"filtered"`
	Items    map[string][]T `json:"items"`
}

type Filter struct {
	Regex     string
	Class     string
	Exception string
	Queue     string
	StartDate time.Time
	EndDate   time.Time
	Filtered  int
}

func ShouldFilterString(f Filter, queue string) bool {
	if f.Regex == "" {
		return false
	}
	matches, _ := regexp.MatchString(f.Regex, queue)
	return !matches
}
