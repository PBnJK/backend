package handler

import (
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/pbnjk/backend/unit-converter/converter"
)

func HandleTemperature(w http.ResponseWriter, r *http.Request) {
	number, err := strconv.ParseFloat(r.FormValue("t-number"), 64)
	if err != nil {
		log.Fatal(err)
	}

	from := r.FormValue("t-from")
	to := r.FormValue("t-to")

	result := converter.ConvertTemperature(number, from, to)

	tmpl := template.Must(template.ParseFiles("site/index.html"))
	if err := tmpl.Execute(w, map[string]interface{}{
		"TNumber": number,
		"TUnit":   getUnitFromID(from),
		"TTo":     getUnitFromID(to),
		"TResult": result,
		"Focus":   "tab-temperature",
	}); err != nil {
		log.Fatal(err)
	}
}
