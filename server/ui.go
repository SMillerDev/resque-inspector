package server

import (
	"embed"
	"html/template"
	"log"
	"net/http"
	"resque-inspector/models"
	"resque-inspector/resque"
)

//go:embed `img/favicon.ico`
var favicon string

//go:embed `template`
var LayoutFs embed.FS

var bootstrap *template.Template
var Dsn string

func getUi(w http.ResponseWriter, r *http.Request) {
	log.Default().Printf("[Web] %s request\n", r.RequestURI)
	page := r.PathValue("page")
	if page == "favicon.ico" {
		w.Header().Set("Content-Type", "img/ico")
		w.Write([]byte(favicon))
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	var err error
	bootstrap, err = template.ParseFS(LayoutFs, "template/*.html")
	if err != nil {
		log.Default().Fatalf("failed to parse bootstrap template: %v", err)
	}

	var templateName string
	switch page {
	case "queues":
		templateName = "queues.html"
	case "workers":
		templateName = "workers.html"
	default:
		templateName = "index.html"
	}
	err = bootstrap.ExecuteTemplate(w, templateName, map[string]interface{}{"Dsn": Dsn, "Page": page, "queues": models.GetQueueList(resque.Filter{}).Items})
	if err != nil {
		//w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}
}
