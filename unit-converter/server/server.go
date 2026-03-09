package server

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/pbnjk/backend/unit-converter/handler"
)

func Run() {
	setupServer()
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func setupServer() {
	setupAllFileHandlers()
	setupAllPostHandlers()
}

// Helper function for setting up all of the needed file handlers
func setupAllFileHandlers() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles("site/index.html"))

		if err := tmpl.Execute(w, nil); err != nil {
			log.Fatal(err)
		}
	})

	setupFileHandler("/style.css", "style.css")
	setupFileHandler("/script.js", "script.js")
}

// Helper function for setting up file handlers
func setupFileHandler(pattern, name string) {
	http.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filepath.Join("site/", name))
	})
}

// Helper function for setting up all of the needed POST request handlers
func setupAllPostHandlers() {
	http.HandleFunc("/length", handler.HandleLength)
	http.HandleFunc("/temperature", handler.HandleTemperature)
	http.HandleFunc("/weight", handler.HandleWeight)
}
