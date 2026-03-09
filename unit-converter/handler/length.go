package handler

import (
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/pbnjk/backend/unit-converter/converter"
)

func HandleLength(w http.ResponseWriter, r *http.Request) {
	number, err := strconv.ParseFloat(r.FormValue("l-number"), 64)
	if err != nil {
		log.Fatal(err)
	}

	from := r.FormValue("l-from")
	to := r.FormValue("l-to")

	result := converter.ConvertLength(number, from, to)

	tmpl := template.Must(template.ParseFiles("site/index.html"))
	if err := tmpl.Execute(w, map[string]interface{}{
		"LNumber": number,
		"LUnit":   getUnitFromID(from),
		"LTo":     getUnitFromID(to),
		"LResult": result,
		"Focus":   "tab-length",
	}); err != nil {
		log.Fatal(err)
	}
}
