package server

import (
	"net/http"
)

func multiplexer() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /collections", renderCollections)
	trace(_https, "handler registered for /collections")
	// /collections/{$} returns 404

	mux.HandleFunc("GET /tags", renderTags)
	trace(_https, "handler registered for /tags")
	// /tags/{$} returns 404

	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		renderNotFound(w, "")
	})
	trace(_https, "handler registered for /")

	mux.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/collections", http.StatusMovedPermanently)
	})
	trace(_https, "handler registered for /{$}")

	return mux
}

func renderCollections(w http.ResponseWriter, r *http.Request) {
	// only collections containing at least 1 item should be rendered
	collections := make(map[string]string)
	for key := range _database.collections {
		collections[key] = _database.declaredCollections[key]
	}
	renderTemplate(w, Page{
		Key:        "collections",
		Title:      "Collections",
		Supertitle: "All",
		Data:       collections,
	})
}

func renderTags(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, Page{
		Title:      "Tags",
		Supertitle: "All",
	})
}

func renderNotFound(w http.ResponseWriter, title string) {
	if title == "" {
		title = "not found"
	}

	renderTemplate(w, Page{
		Key:        "not-found",
		Title:      title,
		Supertitle: "404",
	})
}
