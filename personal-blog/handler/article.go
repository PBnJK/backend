package handler

import (
	"log"
	"net/http"
	"strconv"

	"github.com/pbnjk/backend/personal-blog/article"
)

func HandleArticle(w http.ResponseWriter, r *http.Request) {
	t := createTemplate("templates/article.html")

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		log.Fatal(err)
	}

	post, err := article.GetArticleByID(id)
	if err != nil {
		log.Fatal(err)
	}

	if err = t.Execute(w, post); err != nil {
		log.Fatal(err)
	}
}
