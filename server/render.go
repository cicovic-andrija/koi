package server

import (
	"html/template"
	"net/http"
)

type Page struct {
	Key        string
	Title      string
	Supertitle string
	Data       any
}

func renderTemplate(w http.ResponseWriter, p Page) {
	template, err := template.ParseFiles("data/pagetemplate.html")
	if err != nil {
		trace(_error, "http: failed to parse template: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = template.Execute(w, p)
	if err != nil {
		trace(_error, "http: template: %v", err)
	}
}
