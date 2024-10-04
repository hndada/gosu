package main

import (
	"html/template"
	"log"
	"net/http"
)

func main() {
	// Handler to serve the template
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Parse all template files
		tmpl := template.Must(template.ParseFiles(
			"layout.html",
			"style.html",
			"left-section.html",
			"middle-section.html",
			"right-section.html",
			"option-panel.html",
			"script.html",
		))

		// Execute only the "layout" template, which includes other templates
		err := tmpl.ExecuteTemplate(w, "layout", map[string]interface{}{
			"Title": "Music Chart Selector",
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	log.Println("Server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
