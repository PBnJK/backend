package handler

import (
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/pbnjk/backend/unit-converter/converter"
)

func HandleWeight(w http.ResponseWriter, r *http.Request) {
	number, err := strconv.ParseFloat(r.FormValue("w-number"), 64)
	if err != nil {
		log.Fatal(err)
	}

	from := r.FormValue("w-from")
	to := r.FormValue("w-to")

	result := converter.ConvertWeight(number, from, to)

	tmpl := template.Must(template.ParseFiles("site/index.html"))
	if err := tmpl.Execute(w, map[string]interface{}{
		"WNumber": number,
		"WUnit":   getUnitFromID(from),
		"WTo":     getUnitFromID(to),
		"WResult": result,
		"Focus":   "tab-weight",
	}); err != nil {
		log.Fatal(err)
	}
}
