package server

import (
	"html/template"
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

	mux.HandleFunc("GET /tags/{tag}", renderTag)
	trace(_https, "handler registered for GET /tags/{tag}")

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
		Key:         "catalogue",
		Title:       _database.declaredCollections[key],
		Supertitle:  "Collection",
		DisplayTags: true,
		Data:        catalogue,
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

func renderTag(w http.ResponseWriter, r *http.Request) {
	tag := r.PathValue("tag")
	catalogue := _database.catalogueOfTaggedItems(tag)
	if catalogue == nil {
		renderNotFound(w, "Tag not found.")
		return
	}
	renderTemplate(w, Page{
		Key:        "catalogue",
		Title:      tag,
		Supertitle: "Items tagged with",
		Data:       catalogue,
	})
}

func renderItem(w http.ResponseWriter, r *http.Request) {
	itemID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || itemID < 0 || itemID > _database.lastID() {
		renderNotFound(w, "Item not found.")
		return
	}

	item := _database.singleItem(itemID)
	renderTemplate(w, Page{
		Key:         item.Type + "/item",
		Title:       item.Label,
		Supertitle:  TypeLabel(item.Type),
		DisplayTags: true,
		Data:        item,
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

// Page is a data structure passed to the template rendering engine.
type Page struct {
	Key         string
	Title       string
	Supertitle  string
	DisplayTags bool
	Data        any
}

// <#hardcoded#>
var pageTemplate = template.Must(template.ParseFiles(
	"data/koipond-main.html",
	"data/koipond-style.html",
	"data/koipond-books.html",
))

func renderTemplate(w http.ResponseWriter, p Page) {
	err := pageTemplate.Execute(w, p)
	if err != nil {
		trace(_error, "http: render template: %v", err)
	}
}
