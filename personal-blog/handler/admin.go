package handler

import (
	"log"
	"net/http"
	"strconv"

	"github.com/pbnjk/backend/personal-blog/article"
)

func HandleAdmin(w http.ResponseWriter, r *http.Request) {
	t := createTemplate("templates/admin.html")
	if err := t.Execute(w, article.GetBlog()); err != nil {
		log.Fatal(err)
	}
}

func HandleAdminNew(w http.ResponseWriter, r *http.Request) {
	// Form was not sent, just normal page access
	if r.Method != http.MethodPost {
		http.ServeFile(w, r, "public/new.html")
		return
	}

	article.AddArticle(r.FormValue("title"), r.FormValue("content"))
	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}

func HandleAdminEdit(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		log.Fatal(err)
	}

	// Form was not sent, just normal page access
	if r.Method != http.MethodPost {
		t := createTemplate("templates/edit.html")

		post, err := article.GetArticleByID(id)
		if err != nil {
			log.Fatal(err)
		}

		if err = t.Execute(w, post); err != nil {
			log.Fatal(err)
		}

		return
	}

	if err := article.EditArticleByID(id, r.FormValue("title"), r.FormValue("content")); err != nil {
		log.Fatal(err)
	}

	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}

func HandleAdminDelete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		log.Fatal(err)
	}

	if err := article.DeleteArticleByID(id); err != nil {
		log.Fatal(err)
	}

	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}
