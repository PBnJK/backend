package handler

import (
	"html/template"
	"log"
	"path/filepath"
	"time"
)

func createTemplate(path string) template.Template {
	base := filepath.Base(path)
	funcMap := template.FuncMap{
		"humanizeDate": humanizeDate,
	}

	return *template.Must(
		template.
			New(base).
			Funcs(funcMap).
			ParseFiles(path),
	)
}

func humanizeDate(date string) string {
	t, err := time.Parse("2006-01-02T15:04:05Z", date)
	if err != nil {
		log.Fatal(err)
	}

	return t.Format("Jan 2 2006, 15:04:05")
}
