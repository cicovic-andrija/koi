package server

import (
	"html/template"
	"net/http"
)

var (
	_pageTemplate = template.Must(
		template.
			New("page").
			ParseFiles(
				"data/main.html",
				"data/style-pretty.html",
				"data/style-minimal.html",
				"data/books.html",
				"data/games.html",
				"data/boardgames.html",
			),
	)

	_customizer = &RenderingCustomizer{
		map[string]bool{
			"@enable-repository-link":  true,
			"@enable-pretty-style":     true,
			"@enable-kanji":            true,
			"@enable-list-decorations": true,
		},
	}
)

// HTMLPage is a main wrapper object sent to the template engine when rendering HTML.
// It contains standard elements of an HTML, e.g. Title, as well as a data object
// that needs to be rendered.
type HTMLPage struct {
	Key          string
	Title        string
	Supertitle   string
	ErrorMessage string
	Data         DataObjectInterface
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
