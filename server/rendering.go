package server

import (
	"html/template"
	"net/http"
)

var (
	// @hardcoded
	_pageTemplate = template.Must(
		template.
			New("page").
			ParseFiles(
				"data/main.html",
				"data/style-pretty.html",
				"data/style-simple.html",
				"data/books.html",
				"data/games.html",
			),
	)

	// @hardcoded
	_customizer = &RenderingCustomizer{
		map[string]bool{
			"@enable-repository-link": true,
			"@enable-pretty-style":    true,
			"@enable-kanji":           false,
		},
	}
)

type HTMLPage struct {
	Key          string
	Title        string
	Supertitle   string
	ErrorMessage string
	Data         CommonDataProperties
}

func (p *HTMLPage) Customizer() *RenderingCustomizer {
	return _customizer
}

func render(w http.ResponseWriter, p HTMLPage) {
	err := _pageTemplate.ExecuteTemplate(w, "main.html", &p)
	if err != nil {
		trace(_error, "http: render template: %v", err)
	}
}
