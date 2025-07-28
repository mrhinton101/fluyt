package web

import (
	"fmt"
	"html/template"
	"net/http"
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.
		New("base.html").
		ParseFiles(
			"web/templates/base.html",
			"web/templates/navbar.html",
			"web/templates/index.html",
		)
	if err != nil {
		http.Error(w, fmt.Sprintf("template error: %v", err), 500)
		return
	}

	data := struct {
		Title string
	}{
		Title: "fluyt - Network Diff",
	}

	err = tmpl.ExecuteTemplate(w, "base", data)
	if err != nil {
		http.Error(w, fmt.Sprintf("execution error: %v", err), 500)
	}
}
