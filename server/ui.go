package server

import (
	"crypto/sha256"
	"embed"
	"encoding/hex"
	"html/template"
	"log"
	"net/http"
	"resque-inspector/models"
)

//go:embed `img/favicon.ico`
var favicon string

//go:embed `js`
var jsFs embed.FS

//go:embed `css`
var cssFs embed.FS

//go:embed `template`
var templateFs embed.FS

var bootstrap *template.Template
var Dsn string

func getUi(w http.ResponseWriter, r *http.Request) {
	log.Default().Printf("[Web] %s %s request\n", r.Method, r.RequestURI)
	if r.Method != "GET" {
		returnError(w, http.StatusMethodNotAllowed, map[string]interface{}{})
		return
	}
	page := r.PathValue("page")
	var err error
	switch r.RequestURI {
	case "/":
		w.Header().Set("Location", "/jobs")
		w.WriteHeader(http.StatusPermanentRedirect)
		return
	case "/favicon.ico":
		w.Header().Set("Content-Type", "img/ico")
		_, err = w.Write([]byte(favicon))
		w.WriteHeader(http.StatusOK)
		return
	case "/js/main.js":
		serveJs(r.RequestURI[1:], w, r)
		return
	case "/css/pico.min.css":
		serveCss(r.RequestURI[1:], w, r)
		return
	case "/css/main.css":
		serveCss(r.RequestURI[1:], w, r)
		return
	default:
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		bootstrap, err = template.ParseFS(templateFs, "template/*.html")
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

func serveCss(filename string, w http.ResponseWriter, r *http.Request) {
	serveStatic(filename, "text/css", cssFs, w, r)
}

func serveJs(filename string, w http.ResponseWriter, r *http.Request) {
	serveStatic(filename, "text/javascript", jsFs, w, r)
}

func serveStatic(filename string, contentType string, fs embed.FS, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", contentType+"; charset=utf-8")
	file, err := fs.ReadFile(filename)
	if err != nil {
		log.Default().Printf("failed to parse %s request: %v", contentType, err)
		w.WriteHeader(http.StatusNotFound)
	}
	etag := computeETag(file)
	if match := r.Header.Get("If-None-Match"); match == etag {
		w.WriteHeader(http.StatusNotModified)
		return
	}

	w.Header().Set("ETag", etag)
	w.Header().Set("Cache-Control", "public, max-age=60")
	_, fileErr := w.Write(file)
	if fileErr != nil {
		log.Default().Printf("failed to serve %s request: %v", contentType, err)
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
}

func computeETag(bytes []byte) string {
	h := sha256.New()
	h.Write(bytes)
	return hex.EncodeToString(h.Sum(nil))
}

func getDataSet(page string, selectedQueue string) map[string]interface{} {
	var data = make(map[string]interface{})
	if page == "queues" || page == "workers" {
		data["page"] = page
	} else {
		data["page"] = "jobs"
	}

	data["dsn"] = Dsn
	data["queues"] = models.GetQueueList(models.Filter{}).Items
	if selectedQueue == "" {
		data["selected"] = "NONE"
	} else {
		data["selected"] = selectedQueue
	}

	return data
}
