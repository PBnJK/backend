package server

import (
	"log"
	"net/http"

	"github.com/pbnjk/backend/personal-blog/handler"
)

func Run() {
	mux := setupServer()
	log.Fatal(http.ListenAndServe(":8080", mux))
}

func setupServer() *http.ServeMux {
	mux := http.NewServeMux()
	setupAllHandlers(mux)

	return mux
}

// Helper function for setting up all of the needed handlers
func setupAllHandlers(mux *http.ServeMux) {
	mux.HandleFunc("/", handler.HandleHome)
	mux.HandleFunc("GET /article/{id}", handler.HandleArticle)

	fs := http.FileServer(http.Dir("public"))
	mux.Handle("/public/", http.StripPrefix("/public/", fs))
}
