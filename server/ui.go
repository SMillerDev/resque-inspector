package server

import (
	_ "embed"
	"fmt"
	"html/template"
	"net/http"
	"resque-inspector/models"
	"resque-inspector/result"
)

var LayoutDir string = "server/template"
var bootstrap *template.Template

//go:embed `img/favicon.ico`
var favicon string

var Dsn string

func getUi(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("got %s UI request\n", r.RequestURI)
	page := r.PathValue("page")
	if page == "favicon.ico" {
		w.Header().Set("Content-Type", "img/ico")
		w.Write([]byte(favicon))
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	var templateName string
	switch page {
	case "queues":
		templateName = "queues.gohtml"
	case "workers":
		templateName = "workers.gohtml"
	default:
		templateName = "index.gohtml"
	}

	var err error
	bootstrap, err = template.ParseGlob(LayoutDir + "/*.gohtml")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	err = bootstrap.ExecuteTemplate(w, templateName, map[string]interface{}{"Dsn": Dsn, "Page": page, "queues": models.GetQueueList(result.Filter{}).Items})
	if err != nil {
		//w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}
}
