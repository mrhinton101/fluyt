package webui

import (
	"fmt"
	"html/template"
	"net/http"
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("webui/templates/index.html")
	if err != nil {
		http.Error(w, fmt.Sprintf("template error: %v", err), http.StatusInternalServerError)
		fmt.Printf("template parse error: %v\n", err)
		return
	}

	err = tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, fmt.Sprintf("template execution error: %v", err), http.StatusInternalServerError)
		fmt.Printf("template execution error: %v\n", err)
	}
}
