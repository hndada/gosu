package main

import (
	"embed"
	"html/template"
	"log"
	"net/http"
)

// embed all template files
//
//go:embed *.html
var content embed.FS

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Parse all template files
		tmpl := template.New("layout")
		tmpl.Funcs(template.FuncMap{
			"embed": func(name string) template.HTML {
				// There are two content variables in this scope.
				// Left one is local variable, right one is package variable.
				content, err := content.ReadFile(name)
				if err != nil {
					return template.HTML("Error reading file: " + name)
				}
				return template.HTML(content)
			},
		})
		_, err := tmpl.ParseFS(content, "*.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Execute only the "layout" template, which includes other templates
		err = tmpl.ExecuteTemplate(w, "layout", map[string]interface{}{
			"Title": "Music Chart Selector",
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	log.Println("Server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// func main0() {
// 	// Handler to serve the template
// 	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
// 		tmpl := template.Must(template.ParseFiles(
// 			"layout.html",
// 			"style.html",
// 			"left-section.html",
// 			"middle-section.html",
// 			"right-section.html",
// 			"option-panel.html",
// 			"script.html",
// 		))
// 		// Execute only the "layout" template, which includes other templates
// 		err := tmpl.ExecuteTemplate(w, "layout", map[string]interface{}{
// 			"Title": "Music Chart Selector",
// 		})
// 		if err != nil {
// 			http.Error(w, err.Error(), http.StatusInternalServerError)
// 		}
// 	})
// 	log.Println("Server started at :8080")
// 	log.Fatal(http.ListenAndServe(":8080", nil))
// }
