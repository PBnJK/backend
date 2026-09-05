package handler

import (
	"log"
	"net/http"

	"github.com/pbnjk/backend/personal-blog/article"
)

func HandleHome(w http.ResponseWriter, r *http.Request) {
	t := createTemplate("templates/home.html")
	if err := t.Execute(w, article.GetBlog()); err != nil {
		log.Fatal(err)
	}
}
