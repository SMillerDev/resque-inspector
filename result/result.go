package result

type Result[T any] struct {
	Filter   string `json:"filter"`
	Total    int    `json:"total"`
	Filtered int    `json:"filtered"`
	Items    []T    `json:"items"`
}

type NamedResult[T any] struct {
	Filter   string         `json:"filter"`
	Total    int            `json:"total"`
	Filtered int            `json:"filtered"`
	Items    map[string][]T `json:"items"`
}
