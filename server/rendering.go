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
				"data/books.html",
				"data/games.html",
				"data/boardgames.html",
			),
	)

	_customizer = &RenderingCustomizer{
		map[string]bool{
			"@enable-kanji":            false,
			"@enable-list-decorations": true,
		},
	}

	_fileServer = http.FileServer(http.Dir("data/static"))
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
