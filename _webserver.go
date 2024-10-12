package gosu

import (
	"html/template"
	"log"
	"net/http"

	"github.com/hndada/gosu/scene/selects3"
)

const HostName = "localhost"
const Port = ":5488"

func RunWebServer() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Parse all template files
		tmpl := template.New("layout")
		tmpl.Funcs(template.FuncMap{
			"embed": func(name string) template.HTML {
				content, err := selects3.Templates.ReadFile(name)
				if err != nil {
					return template.HTML("Error reading file: " + name)
				}
				return template.HTML(content)
			},
		})

		_, err := tmpl.ParseFS(selects3.Templates, "*.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Execute only the "layout" template
		err = tmpl.ExecuteTemplate(w, "layout", map[string]interface{}{
			"Title": "Music Chart Selector",
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	log.Println("Server started at", Port)
	log.Fatal(http.ListenAndServe(Port, nil))
}
