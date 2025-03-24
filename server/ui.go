package server

import (
	_ "embed"
	"fmt"
	"html/template"
	"net/http"
)

//go:embed `index.html.tmpl`
var page string
var Dsn string

func GetRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("got %s request\n", r.RequestURI)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	templ, err := template.New("ui").Parse(page)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	err = templ.Execute(w, map[string]interface{}{"Dsn": Dsn})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
