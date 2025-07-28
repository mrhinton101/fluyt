package web

import (
	"fmt"
	"html/template"
	"net/http"
)

func renderTemplate(w http.ResponseWriter, templateName string, data any) {
	tmpl, err := template.
		New("base.html").
		ParseFiles(
			"web/templates/base.html",
			"web/templates/navbar.html",
			fmt.Sprintf("web/templates/%s", templateName),
		)
	if err != nil {
		http.Error(w, fmt.Sprintf("template error: %v", err), 500)
		return
	}

	err = tmpl.ExecuteTemplate(w, "base", data)
	if err != nil {
		http.Error(w, fmt.Sprintf("execution error: %v", err), 500)
	}
}
