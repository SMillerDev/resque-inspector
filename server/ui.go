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

//go:embed `template/main.js`
var js string

//go:embed `template`
var LayoutFs embed.FS

var bootstrap *template.Template
var Dsn string

func getUi(w http.ResponseWriter, r *http.Request) {
	log.Default().Printf("[Web] %s request\n", r.RequestURI)
	page := r.PathValue("page")
	var err error
	switch r.RequestURI {
	case "/favicon.ico":
		w.Header().Set("Content-Type", "img/ico")
		_, err = w.Write([]byte(favicon))
		w.WriteHeader(200)
		return
	case "/js/main.js":
		w.Header().Set("Content-Type", "text/javascript; charset=utf-8")
		_, err = w.Write([]byte(js))
		w.WriteHeader(200)
		return
	default:
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		bootstrap, err = template.ParseFS(LayoutFs, "template/*.html")
	}

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

	selected := r.URL.Query().Get("queue")
	err = bootstrap.ExecuteTemplate(w, templateName, getDataSet(page, selected))
	if err != nil {
		//w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}
}

func getDataSet(page string, selectedQueue string) map[string]interface{} {
	var data = make(map[string]interface{})
	if page == "queues" || page == "workers" {
		data["page"] = page
	} else {
		data["page"] = "jobs"
	}

	data["dsn"] = Dsn
	data["queues"] = models.GetQueueList(resque.Filter{}).Items
	if selectedQueue == "" {
		data["selected"] = "NONE"
	} else {
		data["selected"] = selectedQueue
	}

	return data
}
