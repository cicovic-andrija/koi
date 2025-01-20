package server

import (
	"net/http"
	"strconv"
)

func multiplexer() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /collections", renderCollections) // /collections/{$}: 404
	trace(_https, "handler registered for GET /collections")

	mux.HandleFunc("GET /collections/{key}", renderCollection)
	trace(_https, "handler registered for GET /collections/{key}")

	mux.HandleFunc("GET /tags", renderTags) // /tags/{$}: 404
	trace(_https, "handler registered for GET /tags")

	mux.HandleFunc("GET /items/{id}", renderItem)
	trace(_https, "handler registered for GET /items/{id}")

	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		renderNotFound(w, "Page not found.")
	})
	trace(_https, "handler registered for GET /")

	mux.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/collections", http.StatusMovedPermanently)
	})
	trace(_https, "handler registered for GET /{$}")

	return mux
}

func renderCollections(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, Page{
		Key:        "collections",
		Title:      "Collections",
		Supertitle: "All",
		Data:       _database.collections(),
	})
}

func renderCollection(w http.ResponseWriter, r *http.Request) {
	key := r.PathValue("key")
	catalogue := _database.collectionCatalogue(key)
	if catalogue == nil {
		renderNotFound(w, "Collection not found.")
		return
	}
	renderTemplate(w, Page{
		Key:        "catalogue",
		Title:      _database.declaredCollections[key],
		Supertitle: "Collection",
		Data:       catalogue,
	})
}

func renderTags(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, Page{
		Key:        "tags",
		Title:      "Tags",
		Supertitle: "All",
		Data:       _database.tags(),
	})
}

func renderItem(w http.ResponseWriter, r *http.Request) {
	itemID := convertAndCheckID(r.PathValue("id"), _database.lastID())
	item := _database.singleItem(itemID)
	renderTemplate(w, Page{
		Key:        "item",
		Title:      item.Label,
		Supertitle: item.TypeLabel(),
	})
}

func renderNotFound(w http.ResponseWriter, content string) {
	renderTemplate(w, Page{
		Key:        "not-found",
		Title:      "Not found",
		Supertitle: "404",
		Data:       content,
	})
}

func convertAndCheckID(idStr string, max int) int {
	id, err := strconv.Atoi(idStr)
	if err != nil || id < 1 || id > max {
		return 0
	}
	return id
}
